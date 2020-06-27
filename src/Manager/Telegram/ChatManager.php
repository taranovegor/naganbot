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
use App\Event\Telegram\ChatEvent;
use App\Exception\Common\EntityNotFoundException;
use App\Repository\Telegram\ChatRepository;
use App\Service\Common\EventDispatcher;
use App\Service\Common\Flusher;
use App\ValueObject\Telegram\Chat\Type;
use App\ValueObject\Telegram\User\UserChatStatistic;
use Doctrine\ORM\ORMException;

/**
 * Class ChatManager
 */
class ChatManager
{
    private ChatRepository $repository;

    private Flusher $flusher;

    private EventDispatcher $eventDispatcher;

    /**
     * ChatManager constructor.
     *
     * @param ChatRepository  $repository
     * @param Flusher         $flusher
     * @param EventDispatcher $eventDispatcher
     */
    public function __construct(ChatRepository $repository, Flusher $flusher, EventDispatcher $eventDispatcher)
    {
        $this->repository = $repository;
        $this->flusher = $flusher;
        $this->eventDispatcher = $eventDispatcher;
    }

    /**
     * @param \TelegramBot\Api\Types\Chat $model
     *
     * @return Chat
     *
     * @throws ORMException
     */
    public function create(\TelegramBot\Api\Types\Chat $model): Chat
    {
        $chat = new Chat($model->getId(), new Type($model->getType()));
        $this->repository->add($chat);
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new ChatEvent($chat), ChatEvent::CREATED);

        $this->update($chat, $model);

        return $chat;
    }

    /**
     * @param Chat                        $chat
     * @param \TelegramBot\Api\Types\Chat $model
     *
     * @return Chat
     */
    public function update(Chat $chat, \TelegramBot\Api\Types\Chat $model): Chat
    {
        $chat
            ->setType(new Type($model->getType()))
            ->setTitle($model->getTitle())
            ->setFirstName($model->getFirstName())
            ->setLastName($model->getLastName())
            ->setInviteLink($model->getInviteLink())
        ;
        $this->flusher->flush();

        $this->eventDispatcher->dispatch(new ChatEvent($chat), ChatEvent::UPDATED);

        return $chat;
    }

    /**
     * @param int $id
     *
     * @return Chat
     *
     * @throws EntityNotFoundException
     */
    public function get(int $id): Chat
    {
        $object = $this->repository->find($id);
        if (!$object instanceof Chat) {
            throw new EntityNotFoundException();
        }

        return $object;
    }

    /**
     * @param Chat $chat
     * @param int  $limit
     *
     * @return UserChatStatistic[]
     */
    public function getUserChatStatisticForEachMemberInChat(Chat $chat, int $limit = 10): array
    {
        return array_map(
            fn(array $r) => new UserChatStatistic($r['user'], $r['number_of_wins']),
            $this->repository->numberOfWinsForEachChatMemberInChat($chat, $limit)
        );
    }
}
