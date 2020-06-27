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

namespace App\Service\Telegram;

use App\Constant\Common\Logger\Channel;
use App\Constant\Telegram\ChatType;
use App\Exception\Common\EntityNotFoundException;
use App\Exception\Telegram\AvailableOnlyInGroupException;
use App\Manager\Telegram\ChatManager;
use App\Manager\Telegram\UserManager;
use App\Service\Common\LoggerFactory;
use Doctrine\ORM\ORMException;
use Psr\Log\LoggerInterface;
use TelegramBot\Api\Types\Update;

/**
 * Class TelegramUpdateHandler
 */
class UpdateHandler
{
    private UserManager $userManager;

    private ChatManager $chatManager;

    private LoggerInterface $logger;

    /**
     * TelegramDataSaver constructor.
     *
     * @param UserManager   $userManager
     * @param ChatManager   $chatManager
     * @param LoggerFactory $loggerFactory
     */
    public function __construct(UserManager $userManager, ChatManager $chatManager, LoggerFactory $loggerFactory)
    {
        $this->userManager = $userManager;
        $this->chatManager = $chatManager;
        $this->logger = $loggerFactory->create(Channel::UPDATE_HANDLER);
    }

    /**
     * @param Update $update
     *
     * @throws AvailableOnlyInGroupException
     * @throws ORMException
     */
    public function handle(Update $update): void
    {
        $message = $update->getMessage();
        if (null === $message) {
            return;
        }

        if (!in_array($message->getChat()->getType(), [
            ChatType::GROUP,
            ChatType::SUPERGROUP,
        ], true)) {
            throw new AvailableOnlyInGroupException();
        }

        try {
            $this->userManager->update(
                $this->userManager->get($message->getFrom()->getId()),
                $message->getFrom()
            );
        } catch (EntityNotFoundException $e) {
            $this->userManager->create($message->getFrom());
        }

        try {
            $this->chatManager->update(
                $this->chatManager->get($message->getChat()->getId()),
                $message->getChat()
            );
        } catch (EntityNotFoundException $e) {
            $this->chatManager->create($message->getChat());
        }
    }
}
