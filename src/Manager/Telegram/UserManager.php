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

namespace App\Manager\Telegram;

use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Event\Telegram\UserEvent;
use App\Exception\Common\EntityNotFoundException;
use App\Repository\Telegram\UserRepository;
use App\Service\Common\EventDispatcher;
use App\Service\Common\Flusher;
use App\ValueObject\Telegram\User\LanguageCode;
use App\ValueObject\Telegram\User\UserChatStatistic;
use App\ValueObject\Telegram\User\Username;
use Doctrine\ORM\NonUniqueResultException;
use Doctrine\ORM\NoResultException;
use Doctrine\ORM\ORMException;

/**
 * Class UserManager
 */
class UserManager
{
    private UserRepository $repository;

    private Flusher $flusher;

    private EventDispatcher $eventDispatcher;

    /**
     * UserManager constructor.
     *
     * @param UserRepository  $repository
     * @param Flusher         $flusher
     * @param EventDispatcher $eventDispatcher
     */
    public function __construct(UserRepository $repository, Flusher $flusher, EventDispatcher $eventDispatcher)
    {
        $this->repository = $repository;
        $this->flusher = $flusher;
        $this->eventDispatcher = $eventDispatcher;
    }

    /**
     * @param \TelegramBot\Api\Types\User $model
     *
     * @return User
     *
     * @throws ORMException
     */
    public function create(\TelegramBot\Api\Types\User $model): User
    {
        $user = new User($model->getId(), $model->getFirstName(), $model->isBot());
        $this->repository->add($user);
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new UserEvent($user), UserEvent::CREATED);

        $this->update($user, $model);

        return $user;
    }

    /**
     * @param User                        $user
     * @param \TelegramBot\Api\Types\User $model
     *
     * @return User
     */
    public function update(User $user, \TelegramBot\Api\Types\User $model): User
    {
        $user
            ->setFirstName($model->getFirstName())
            ->setLastName($model->getLastName())
            ->setUsername(new Username($model->getUsername()))
            ->setIsBot($model->isBot())
            ->setLanguageCode(new LanguageCode($model->getLanguageCode()))
        ;
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new UserEvent($user), UserEvent::UPDATED);

        return $user;
    }

    /**
     * @param int $id
     *
     * @return User
     *
     * @throws EntityNotFoundException
     */
    public function get(int $id): User
    {
        $user = $this->repository->find($id);
        if (!$user instanceof User) {
            throw new EntityNotFoundException();
        }

        return $user;
    }

    /**
     * @param User $user
     * @param Chat $chat
     *
     * @return UserChatStatistic
     *
     * @throws NoResultException
     * @throws NonUniqueResultException
     */
    public function getUserChatStatisticForUserInChat(User $user, Chat $chat): UserChatStatistic
    {
        return new UserChatStatistic(
            $user,
            $this->repository->numberOfWinsForUserInChat($user, $chat)
        );
    }
}
