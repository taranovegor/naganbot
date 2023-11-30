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
use App\Entity\Telegram\User;
use App\Event\Game\GunslingerEvent;
use App\Exception\Common\EntityNotFoundException;
use App\Exception\Game\AlreadyJoinedToGameException;
use App\Exception\Game\AlreadyPlayedException;
use App\Repository\Game\GunslingerRepository;
use App\Service\Common\EventDispatcher;
use App\Service\Common\Flusher;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\ORMException;

/**
 * Class GunslingerManager
 */
class GunslingerManager
{
    private GunslingerRepository $repository;

    private Flusher $flusher;

    private EventDispatcher $eventDispatcher;

    /**
     * GunslingerManager constructor.
     *
     * @param GunslingerRepository $repository
     * @param Flusher              $flusher
     * @param EventDispatcher      $eventDispatcher
     */
    public function __construct(GunslingerRepository $repository, Flusher $flusher, EventDispatcher $eventDispatcher)
    {
        $this->repository = $repository;
        $this->flusher = $flusher;
        $this->eventDispatcher = $eventDispatcher;
    }

    /**
     * @param Game $game
     * @param User $user
     *
     * @return Gunslinger
     *
     * @throws AlreadyJoinedToGameException
     * @throws AlreadyPlayedException
     * @throws ORMException
     */
    public function create(Game $game, User $user): Gunslinger
    {
        if ($game->isPlayed() && $game->isCreatedToday()) {
            throw new AlreadyPlayedException();
        }

        try {
            $this->getByUserInGame($user, $game);

            throw new AlreadyJoinedToGameException();
        } catch (EntityNotFoundException $e) {
        }

        $gunslinger = new Gunslinger($game, $user);
        $this->repository->add($gunslinger);
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new GunslingerEvent($gunslinger), GunslingerEvent::JOINED_TO_GAME);

        return $gunslinger;
    }

    /**
     * @param Game $game
     *
     * @return Collection
     */
    public function getByGame(Game $game): Collection
    {
        return new ArrayCollection($this->repository->findByGame($game));
    }

    /**
     * @param User $user
     * @param Game $game
     *
     * @return Gunslinger
     *
     * @throws EntityNotFoundException
     */
    public function getByUserInGame(User $user, Game $game): Gunslinger
    {
        $gunslinger = $this->repository->findByUserInGame($user, $game);
        if (!$gunslinger instanceof Gunslinger) {
            throw new EntityNotFoundException();
        }

        return $gunslinger;
    }

    /**
     * @param Gunslinger $gunslinger
     *
     * @return Gunslinger
     */
    public function shot(Gunslinger $gunslinger): Gunslinger
    {
        $gunslinger->shot();
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new GunslingerEvent($gunslinger), GunslingerEvent::DIED);

        return $gunslinger;
    }
}
