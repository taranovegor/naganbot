<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Service;

use App\Entity\Game;
use App\Entity\Gunslinger;
use App\Entity\ShotHimself;
use App\Event\GameEvent;
use App\Event\GunslingerJoinedEvent;
use App\Event\GunslingerShotHimselfEvent;
use App\Exception\Game\AlreadyRegisteredInGameException;
use App\Exception\EntityNotFoundException;
use App\Exception\Game\FailedToScrollDrumException;
use App\Exception\Game\GameIsAlreadyCreatedException;
use App\Exception\Game\ShotDeadNotFoundException;
use App\Exception\Game\GameIsAlreadyPlayedException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Game\NotFoundActiveGameException;
use App\Repository\GameRepository;
use App\Repository\GunslingerRepository;
use App\Repository\ShotHimselfRepository;
use App\Repository\Telegram\ChatRepository;
use App\Repository\Telegram\UserRepository;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\ORMException;
use Exception;
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
    private GameRepository $gameTableRepository;

    /**
     * @var GunslingerRepository
     */
    private GunslingerRepository $gunslingerRepository;

    /**
     * @var ShotHimselfRepository
     */
    private ShotHimselfRepository $shotHimselfRepository;

    /**
     * @var Flusher
     */
    private Flusher $flusher;

    /**
     * @var EventDispatcher
     */
    private EventDispatcher $eventDispatcher;

    /**
     * GameManager constructor.
     *
     * @param ChatRepository        $chatRepository
     * @param UserRepository        $userRepository
     * @param GameRepository        $gameTableRepository
     * @param GunslingerRepository  $gunslingerRepository
     * @param ShotHimselfRepository $shotHimselfRepository
     * @param Flusher               $flusher
     * @param EventDispatcher       $eventDispatcher
     */
    public function __construct(ChatRepository $chatRepository, UserRepository $userRepository, GameRepository $gameTableRepository, GunslingerRepository $gunslingerRepository, ShotHimselfRepository $shotHimselfRepository, Flusher $flusher, EventDispatcher $eventDispatcher)
    {
        $this->chatRepository = $chatRepository;
        $this->userRepository = $userRepository;
        $this->gameTableRepository = $gameTableRepository;
        $this->gunslingerRepository = $gunslingerRepository;
        $this->shotHimselfRepository = $shotHimselfRepository;
        $this->flusher = $flusher;
        $this->eventDispatcher = $eventDispatcher;
    }

    /**
     * @param Chat $chat
     * @param User $inviter
     *
     * @return Game
     *
     * @throws GameIsAlreadyCreatedException
     * @throws GameIsAlreadyPlayedException
     * @throws ORMException
     * @throws EntityNotFoundException
     */
    public function create(Chat $chat, User $inviter): Game
    {
        $chat = $this->chatRepository->get($chat->getId());
        $user = $this->userRepository->get($inviter->getId());

        try {
            $game = $this->gameTableRepository->getByChat($chat);
            if (false === $game->isPlayed()) {
                throw new GameIsAlreadyCreatedException();
            }
            if($game->isCreatedToday()) {
                throw new GameIsAlreadyPlayedException();
            }
        } catch (EntityNotFoundException $e) {
            // ignore
        }

        $game = new Game($chat, $user);
        $gunslinger = new Gunslinger($game, $user);
        $this->gameTableRepository->add($game);
        $this->gunslingerRepository->add($gunslinger);

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
     * @throws NotFoundActiveGameException
     * @throws GameIsAlreadyPlayedException
     * @throws AlreadyRegisteredInGameException
     * @throws EntityNotFoundException
     * @throws ORMException
     */
    public function join(Chat $chat, User $entering): Gunslinger
    {
        $chat = $this->chatRepository->get($chat->getId());
        try {
            $game = $this->gameTableRepository->getActiveOrCreatedTodayByChat($chat);
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
        $this->gunslingerRepository->add($gunslinger);

        $this->flusher->flush();

        $this->eventDispatcher->dispatch(
            new GunslingerJoinedEvent($game, $gunslinger),
            GunslingerJoinedEvent::EVENT
        );

        return $gunslinger;
    }

    /**
     * @param Chat $chat
     *
     * @return ShotHimself
     *
     * @throws EntityNotFoundException
     * @throws NotEnoughGunslingersException
     * @throws ORMException
     * @throws ShotDeadNotFoundException
     * @throws Exception
     */
    public function play(Chat $chat): ShotHimself
    {
        $chat = $this->chatRepository->get($chat->getId());
        $game = $this->gameTableRepository->getByChat($chat);
        $gunslingers = $game->getGunslingers();

        if ($gunslingers->count() < 6) {
            throw new NotEnoughGunslingersException();
        }

        $drum = [0, 0, 0, 0, 0, 1];
        if (false === shuffle($drum)) {
            throw new FailedToScrollDrumException();
        }

        $chamber = array_search(1, $drum, true);
        $gunslinger = $gunslingers->get($chamber);
        if (!$gunslinger instanceof Gunslinger) {
            throw new ShotDeadNotFoundException();
        }

        $shotHimself = new ShotHimself($gunslinger);
        $this->shotHimselfRepository->add($shotHimself);
        $game->setAsPlayed();

        $this->flusher->flush();

        $this->eventDispatcher->dispatch(
            new GunslingerShotHimselfEvent($game, $shotHimself),
            GunslingerShotHimselfEvent::EVENT
        );

        return $shotHimself;
    }

    /**
     * @param Chat $chat
     *
     * @return Gunslinger[]|Collection
     * @throws EntityNotFoundException
     */
    public function joined(Chat $chat): Collection
    {
        $chat = $this->chatRepository->get($chat->getId());
        $gameTable = $this->gameTableRepository->getByChat($chat);

        return $gameTable->getGunslingers();
    }
}
