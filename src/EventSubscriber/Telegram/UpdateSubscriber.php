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

namespace App\EventSubscriber\Telegram;

use App\Exception\Telegram\AvailableOnlyInGroupException;
use App\Service\Telegram\UpdateHandler;
use BoShurik\TelegramBotBundle\Event\UpdateEvent;
use Doctrine\ORM\ORMException;
use Psr\Log\LoggerInterface;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;

/**
 * Class UpdateSubscriber
 */
class UpdateSubscriber implements EventSubscriberInterface
{
    private UpdateHandler $updateHandler;

    private LoggerInterface $logger;

    /**
     * CommandSubscriber constructor.
     *
     * @param UpdateHandler   $dataSaver
     * @param LoggerInterface $logger
     */
    public function __construct(UpdateHandler $dataSaver, LoggerInterface $logger)
    {
        $this->updateHandler = $dataSaver;
        $this->logger = $logger;
    }

    /**
     * @inheritDoc
     */
    public static function getSubscribedEvents()
    {
        return [
            UpdateEvent::class => [
                ['handleUpdate', 255],
            ],
        ];
    }

    /**
     * @param UpdateEvent $event
     *
     * @throws ORMException
     * @throws AvailableOnlyInGroupException
     */
    public function handleUpdate(UpdateEvent $event)
    {
        $this->updateHandler->handle($event->getUpdate());
    }
}
