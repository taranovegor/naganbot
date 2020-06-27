<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Manager;

use App\Entity\Game\Game;
use App\Entity\Game\Gunslinger;
use App\Event\GameEvent;
use App\Event\GunslingerEvent;
use App\Exception\Game\AlreadyRegisteredInGameException;
use App\Exception\EntityNotFoundException;
use App\Exception\Game\FailedToScrollDrumException;
use App\Exception\Game\GameIsAlreadyCreatedException;
use App\Exception\Game\ShotDeadNotFoundException;
use App\Exception\Game\GameIsAlreadyPlayedException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Game\NotFoundActiveGameException;
use App\Repository\Game\GameRepository;
use App\Repository\Game\GunslingerRepository;
use App\Repository\Telegram\ChatRepository;
use App\Repository\Telegram\UserRepository;
use App\Service\EventDispatcher;
use App\Service\Flusher;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\ORMException;
use TelegramBot\Api\Types\Chat;
use TelegramBot\Api\Types\User;

/**
 * Class GameManager
 */
class GameManager
{
    /**
     * @var ChatRepository
     */
    private ChatRepository $chatRepository;

    /**
     * @var UserRepository
     */
    private UserRepository $userRepository;

    /**
     * @var GameRepository
     */
    private GameRepository $repository;

    /**
     * @var GunslingerRepository
     */
    private GunslingerRepository $gunslingerRepository;

    /**
     * @var Flusher
     */
    private Flusher $flusher;

    /**
     * @var EventDispatcher
     */
    private EventDispatcher $eventDispatcher;

    /**
     * @var int
     */
    private int $numberOfPlayers;

    /**
     * GameManager constructor.
     *
     * @param ChatRepository       $chatRepository
     * @param UserRepository       $userRepository
     * @param GameRepository       $gameTableRepository
     * @param GunslingerRepository $gunslingerRepository
     * @param Flusher              $flusher
     * @param EventDispatcher      $eventDispatcher
     * @param int                  $numberOfPlayers
     */
    public function __construct(ChatRepository $chatRepository, UserRepository $userRepository, GameRepository $gameTableRepository, GunslingerRepository $gunslingerRepository, Flusher $flusher, EventDispatcher $eventDispatcher, int $numberOfPlayers)
    {
        $this->chatRepository = $chatRepository;
        $this->userRepository = $userRepository;
        $this->repository = $gameTableRepository;
        $this->gunslingerRepository = $gunslingerRepository;
        $this->flusher = $flusher;
        $this->eventDispatcher = $eventDispatcher;
        $this->numberOfPlayers = $numberOfPlayers;
    }

    /**
     * @param \App\Entity\Telegram\Chat $chat
     *
     * @return Game
     *
     * @throws EntityNotFoundException
     */
    public function getLatestByChat(\App\Entity\Telegram\Chat $chat): Game
    {
        $object = $this->repository->findLatestByChat($chat);
        if (!$object instanceof Game) {
            throw new EntityNotFoundException();
        }

        return $object;
    }

    /**
     * @param Chat $chat
     * @param User $inviting
     *
     * @return Game
     *
     * @throws GameIsAlreadyCreatedException
     * @throws GameIsAlreadyPlayedException
     * @throws EntityNotFoundException
     */
    public function create(Chat $chat, User $inviting): Game
    {
        $chat = $this->chatRepository->get($chat->getId());
        $user = $this->userRepository->get($inviting->getId());

        try {
            $game = $this->repository->getByChat($chat);
            if (false === $game->isPlayed()) {
                throw new GameIsAlreadyCreatedException();
            }
            if ($game->isCreatedToday()) {
                throw new GameIsAlreadyPlayedException();
            }
        } catch (EntityNotFoundException $e) {
            // ignore
        }

        $game = new Game($chat, $user);
        $gunslinger = new Gunslinger($game, $user);
        $game->addGunslinger($gunslinger);

        $this->flusher->flush();

        $this->eventDispatcher->dispatch(
            new GameEvent($game),
            GameEvent::CREATED
        );

        return $game;
    }

    /**
     * @param Chat $chat
     * @param User $entering
     *
     * @return Gunslinger
     *
     * @throws AlreadyRegisteredInGameException
     * @throws EntityNotFoundException
     * @throws GameIsAlreadyPlayedException
     * @throws NotFoundActiveGameException
     */
    public function join(Chat $chat, User $entering): Gunslinger
    {
        $chat = $this->chatRepository->get($chat->getId());
        try {
            $game = $this->repository->getActiveOrCreatedTodayByChat($chat);
        } catch (EntityNotFoundException $e) {
            throw new NotFoundActiveGameException();
        }
        if ($game->getPlayedAt() && $game->isCreatedToday()) {
            throw new GameIsAlreadyPlayedException();
        }

        $user = $this->userRepository->get($entering->getId());
        if ($game->getGunslingerByUser($user) instanceof Gunslinger) {
            throw new AlreadyRegisteredInGameException();
        }
        $gunslinger = new Gunslinger($game, $user);
        $game->addGunslinger($gunslinger);

        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new GunslingerEvent($gunslinger), GunslingerEvent::JOINED);

        return $gunslinger;
    }

    /**
     * @param Chat $chat
     *
     * @return Gunslinger
     *
     * @deprecated use \App\Manager\GameManager::playGame instead
     *
     * @throws EntityNotFoundException
     * @throws FailedToScrollDrumException
     * @throws NotEnoughGunslingersException
     * @throws ShotDeadNotFoundException
     * @throws GameIsAlreadyPlayedException
     */
    public function play(Chat $chat): Gunslinger
    {
        $chat = $this->chatRepository->get($chat->getId());
        $game = $this->getLatestByChat($chat);

        return $this->playGame($game);
    }

    /**
     * @param Game $game
     *
     * @return Gunslinger
     *
     * @throws FailedToScrollDrumException
     * @throws GameIsAlreadyPlayedException
     * @throws NotEnoughGunslingersException
     * @throws ShotDeadNotFoundException
     */
    public function playGame(Game $game): Gunslinger
    {
        $gunslingers = $game->getGunslingers();

        if ($game->isPlayed()) {
            throw new GameIsAlreadyPlayedException();
        }

        if (false === $this->isEnoughPlayers($game)) {
            throw new NotEnoughGunslingersException();
        }

        $drum = [0, 0, 0, 0, 0, 1];
        if (false === shuffle($drum)) {
            throw new FailedToScrollDrumException();
        }

        $chamber = (int) array_search(1, $drum, true);
        $gunslinger = $gunslingers->get($chamber);
        if (!$gunslinger instanceof Gunslinger) {
            throw new ShotDeadNotFoundException();
        }

        $gunslinger->shot();
        $game->setAsPlayed();

        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new GameEvent($game), GameEvent::PLAYED);
        $this->eventDispatcher->dispatch(new GunslingerEvent($gunslinger), GunslingerEvent::SHOT_HIMSELF);

        return $gunslinger;
    }

    /**
     * @param Chat $chat
     *
     * @return Gunslinger[]|Collection
     *
     * @throws EntityNotFoundException
     */
    public function joined(Chat $chat): Collection
    {
        $chat = $this->chatRepository->get($chat->getId());
        $gameTable = $this->repository->getByChat($chat);

        return $gameTable->getGunslingers();
    }

    /**
     * @param Game $game
     *
     * @return bool
     */
    public function isEnoughPlayers(Game $game): bool
    {
        return $game->getGunslingers()->count() >= $this->numberOfPlayers;
    }
}
