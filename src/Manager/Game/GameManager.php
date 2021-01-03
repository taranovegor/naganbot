<?php
/**
 * Copyright (C) 14.08.20 Egor Taranov
 * This file is part of Nagan bot <https://github.com/taranovegor/nagan-bot>.
 *
 * Nagan bot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Nagan bot is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Nagan bot.  If not, see <http://www.gnu.org/licenses/>.
 */

namespace App\Manager\Game;

use App\Entity\Game\Game;
use App\Entity\Game\Gunslinger;
use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Event\Game\GameEvent;
use App\Exception\Common\EntityNotFoundException;
use App\Exception\Game\AlreadyCreatedException;
use App\Exception\Game\AlreadyJoinedToGameException;
use App\Exception\Game\AlreadyPlayedException;
use App\Exception\Game\FailedToShuffleArrayException;
use App\Exception\Game\GunslingerNotFoundException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Repository\Game\GameRepository;
use App\Service\Common\EventDispatcher;
use App\Service\Common\Flusher;
use App\Service\Game\NuclearBulletChecker;
use Doctrine\ORM\ORMException;

/**
 * Class GameManager
 */
class GameManager
{
    private GameRepository $repository;

    private Flusher $flusher;

    private EventDispatcher $eventDispatcher;

    private GunslingerManager $gunslingerManager;

    private NuclearBulletChecker $nuclearBulletChecker;

    private int $numberOfPlayers;

    /**
     * GameManager constructor.
     *
     * @param GameRepository       $gameTableRepository
     * @param Flusher              $flusher
     * @param EventDispatcher      $eventDispatcher
     * @param GunslingerManager    $gunslingerManager
     * @param NuclearBulletChecker $nuclearBulletChecker
     * @param int                  $numberOfPlayers
     */
    public function __construct(GameRepository $gameTableRepository, Flusher $flusher, EventDispatcher $eventDispatcher, GunslingerManager $gunslingerManager, NuclearBulletChecker $nuclearBulletChecker, int $numberOfPlayers)
    {
        $this->repository = $gameTableRepository;
        $this->flusher = $flusher;
        $this->eventDispatcher = $eventDispatcher;
        $this->gunslingerManager = $gunslingerManager;
        $this->nuclearBulletChecker = $nuclearBulletChecker;
        $this->numberOfPlayers = $numberOfPlayers;
    }

    /**
     * @param Chat $chat
     * @param User $owner
     *
     * @return Game
     *
     * @throws AlreadyCreatedException
     * @throws AlreadyJoinedToGameException
     * @throws AlreadyPlayedException
     * @throws ORMException
     */
    public function create(Chat $chat, User $owner): Game
    {
        try {
            $game = $this->getLatestByChat($chat);

            if (!$game->isPlayed()) {
                throw new AlreadyCreatedException();
            }

            if ($game->isCreatedToday()) {
                throw new AlreadyPlayedException();
            }
        } catch (EntityNotFoundException $e) {
        }

        $game = new Game($chat, $owner);
        $this->repository->add($game);
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new GameEvent($game), GameEvent::CREATED);

        $this->gunslingerManager->create($game, $owner);

        return $game;
    }

    /**
     * @param Chat $chat
     *
     * @return Game
     *
     * @throws EntityNotFoundException
     */
    public function getLatestByChat(Chat $chat): Game
    {
        $game = $this->repository->findLatestByChat($chat);
        if (!$game instanceof Game) {
            throw new EntityNotFoundException();
        }

        return $game;
    }

    /**
     * @param Chat $chat
     *
     * @return Game
     *
     * @throws EntityNotFoundException
     */
    public function getLatestActiveByChat(Chat $chat): Game
    {
        $game = $this->repository->findLatestActiveByChat($chat);
        if (!$game instanceof Game) {
            throw new EntityNotFoundException();
        }

        return $game;
    }

    /**
     * @param Chat $chat
     *
     * @return Game
     *
     * @throws EntityNotFoundException
     */
    public function getActiveOrCreatedTodayByChat(Chat $chat): Game
    {
        $game = $this->repository->findActiveOrCreatedTodayByChat($chat);
        if (!$game instanceof Game) {
            throw new EntityNotFoundException();
        }

        return $game;
    }

    /**
     * @param Game $game
     *
     * @return Game
     *
     * @throws AlreadyPlayedException
     * @throws FailedToShuffleArrayException
     * @throws GunslingerNotFoundException
     * @throws NotEnoughGunslingersException
     */
    public function play(Game $game): Game
    {
        if ($game->isPlayed()) {
            throw new AlreadyPlayedException();
        }
        if (false === $this->isEnoughPlayers($game)) {
            throw new NotEnoughGunslingersException();
        }

        $this->eventDispatcher->dispatch(new GameEvent($game), GameEvent::READY_TO_PLAY);

        if ($this->nuclearBulletChecker->isNuclearBullet($game)) {
            $game = $this->playWithNuclearBullet($game);
        } else {
            $game = $this->playWithRegularBullet($game);
        }

        $game->markAsPlayed();
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new GameEvent($game), GameEvent::PLAYED);

        return $game;
    }

    /**
     * @param Game $game
     *
     * @return bool
     */
    public function isEnoughPlayers(Game $game): bool
    {
        return $this->gunslingerManager->getByGame($game)->count() >= $this->numberOfPlayers;
    }

    /**
     * @param Game $game
     *
     * @return Game
     *
     * @throws FailedToShuffleArrayException
     * @throws GunslingerNotFoundException
     */
    private function playWithRegularBullet(Game $game): Game
    {
        $gunslingers = $this->gunslingerManager->getByGame($game);

        $drum = array_fill(0, $gunslingers->count(), 0);
        $drum[rand(0, $gunslingers->count() - 1)] = 1;
        if (false === shuffle($drum)) {
            throw new FailedToShuffleArrayException();
        }

        $chamber = (int) array_search(1, $drum, true);
        $gunslinger = $gunslingers->get($chamber);
        if (!$gunslinger instanceof Gunslinger) {
            throw new GunslingerNotFoundException();
        }

        $this->gunslingerManager->shot($gunslinger);

        return $game;
    }

    /**
     * @param Game $game
     *
     * @return Game
     */
    private function playWithNuclearBullet(Game $game): Game
    {
        $game->markAsPlayedWithNuclearBullet();

        $gunslingers = $this->gunslingerManager->getByGame($game);
        foreach ($gunslingers as $gunslinger) {
            $this->gunslingerManager->shot($gunslinger);
        }

        return $game;
    }
}
