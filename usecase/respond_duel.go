package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/taranovegor/naganbot/domain"
	"github.com/taranovegor/naganbot/drand"
	"github.com/taranovegor/naganbot/service"
)

var (
	ErrDuelNotOpponent = errors.New("only the challenged user can respond")
	ErrDuelClosed      = errors.New("duel is already closed")
	ErrDuelExpired     = errors.New("duel has expired")
)

type DuelOutcome struct {
	Duel     *domain.Duel
	Loser    domain.User
	ProofURL string
}

type RespondDuelUseCase struct {
	locker service.Locker
	duels  domain.DuelRepository
	drand  *drand.Client
}

func NewRespondDuelUseCase(
	locker service.Locker,
	duels domain.DuelRepository,
	drand *drand.Client,
) *RespondDuelUseCase {
	return &RespondDuelUseCase{locker: locker, duels: duels, drand: drand}
}

func duelLockKey(duelID uuid.UUID) string {
	return fmt.Sprintf("respond-duel-%s", duelID.String())
}

// Accept plays the duel. Only one concurrent call can store the outcome.
func (uc *RespondDuelUseCase) Accept(ctx context.Context, duelID uuid.UUID, userID int64, now time.Time) (*DuelOutcome, error) {
	release, ok := uc.locker.TryLock(duelLockKey(duelID))
	if !ok {
		return nil, service.ErrLockNotAcquired
	}
	defer release()

	duel, err := uc.open(duelID, userID, now)
	if err != nil {
		return nil, err
	}

	loserID, proofURL, err := uc.shoot(ctx, duel)
	if err != nil {
		return nil, err
	}

	duel.MarkAsPlayed(loserID, proofURL)
	claimed, err := uc.duels.Resolve(duel)
	if err != nil {
		return nil, err
	}
	if !claimed {
		return nil, ErrDuelClosed
	}

	loser := duel.Opponent
	if loserID == duel.ChallengerID {
		loser = duel.Challenger
	}

	return &DuelOutcome{Duel: duel, Loser: loser, ProofURL: proofURL}, nil
}

func (uc *RespondDuelUseCase) Decline(duelID uuid.UUID, userID int64, now time.Time) (*domain.Duel, error) {
	release, ok := uc.locker.TryLock(duelLockKey(duelID))
	if !ok {
		return nil, service.ErrLockNotAcquired
	}
	defer release()

	duel, err := uc.open(duelID, userID, now)
	if err != nil {
		return nil, err
	}

	duel.MarkAsDeclined()
	claimed, err := uc.duels.Resolve(duel)
	if err != nil {
		return nil, err
	}
	if !claimed {
		return nil, ErrDuelClosed
	}

	return duel, nil
}

func (uc *RespondDuelUseCase) open(duelID uuid.UUID, userID int64, now time.Time) (*domain.Duel, error) {
	duel, err := uc.duels.GetByID(duelID)
	if err != nil {
		return nil, err
	}
	if duel.IsClosed() {
		return nil, ErrDuelClosed
	}
	if userID != duel.OpponentID {
		return nil, ErrDuelNotOpponent
	}
	if duel.IsExpired(now) {
		return nil, ErrDuelExpired
	}

	return duel, nil
}

func (uc *RespondDuelUseCase) shoot(ctx context.Context, duel *domain.Duel) (loserID int64, proofURL string, err error) {
	beacon, err := uc.drand.GetLatest(ctx)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get drand beacon: %w", err)
	}

	randomness, err := hex.DecodeString(beacon.Randomness)
	if err != nil {
		return 0, "", fmt.Errorf("failed to decode randomness: %w", err)
	}

	hash := sha256.Sum256(append(randomness, []byte(duel.ID.String())...))
	if binary.BigEndian.Uint64(hash[4:12])%2 == 0 {
		loserID = duel.ChallengerID
	} else {
		loserID = duel.OpponentID
	}

	return loserID, uc.drand.ProofURL(beacon.Round, duel.ID), nil
}
