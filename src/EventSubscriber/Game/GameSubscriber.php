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

namespace App\EventSubscriber\Game;

use App\Constant\Common\Logger\Channel;
use App\Constant\Telegram\ParseMode;
use App\Event\Game\GameEvent;
use App\Event\Game\GunslingerEvent;
use App\Exception\Game\AlreadyPlayedException;
use App\Exception\Game\FailedToShuffleArrayException;
use App\Exception\Game\GunslingerNotFoundException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Translation\MessagePatternIsInvalidException;
use App\Manager\Game\GameManager;
use App\Service\Common\LoggerFactory;
use App\Service\MessageBuilder\Game\GameMessageBuilder;
use Psr\Log\LoggerInterface;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Exception;
use TelegramBot\Api\InvalidArgumentException;
use Twig\Error\LoaderError;
use Twig\Error\RuntimeError;
use Twig\Error\SyntaxError;

/**
 * Class GameSubscriber
 */
class GameSubscriber implements EventSubscriberInterface
{
    private GameMessageBuilder $messageBuilder;

    private BotApi $api;

    private LoggerInterface $logger;

    private GameManager $manager;

    /**
     * GameSubscriber constructor.
     *
     * @param GameMessageBuilder $messageBuilder
     * @param BotApi             $api
     * @param LoggerFactory      $loggerFactory
     * @param GameManager        $manager
     */
    public function __construct(GameMessageBuilder $messageBuilder, BotApi $api, LoggerFactory $loggerFactory, GameManager $manager)
    {
        $this->messageBuilder = $messageBuilder;
        $this->api = $api;
        $this->logger = $loggerFactory->create(Channel::GAME);
        $this->manager = $manager;
    }

    /**
     * @inheritDoc
     */
    public static function getSubscribedEvents()
    {
        return [
            GameEvent::CREATED => [
                ['sendMessageWhenCreated', 0],
            ],
            GunslingerEvent::JOINED_TO_GAME => [
                ['checkPossibilityToPlay', 0],
            ],
            GameEvent::READY_TO_PLAY => [
                ['sendMessageWhenReadyToPlay', 0],
            ],
            GameEvent::PLAYED => [
                ['sendMessageWhenPlayedWithNuclearBullet', 0],
            ],
        ];
    }

    /**
     * @param GameEvent $event
     *
     * @throws Exception
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function sendMessageWhenCreated(GameEvent $event): void
    {
        $this->api->sendMessage(
            $event->getGame()->getChat()->getId(),
            $this->messageBuilder->buildCreate(),
            ParseMode::MARKDOWN
        );

        $this->logger->info('event_subscriber.game.game.send_message_when_created', [
            'game.id' => $event->getGame()->getId()->toString(),
        ]);
    }

    /**
     * @param GunslingerEvent $event
     *
     * @throws AlreadyPlayedException
     * @throws FailedToShuffleArrayException
     * @throws GunslingerNotFoundException
     * @throws NotEnoughGunslingersException
     */
    public function checkPossibilityToPlay(GunslingerEvent $event): void
    {
        $game = $event->getGunslinger()->getGame();

        if (!$this->manager->isEnoughPlayers($game)) {
            $this->logger->info('event_subscriber.game.game.check_possibility_to_play.not_enough_players', [
                'game.id' => $event->getGunslinger()->getGame()->getId()->toString(),
            ]);

            return;
        }

        $this->manager->play($game);

        $this->logger->info('event_subscriber.game.game.check_possibility_to_play', [
            'game.id' => $event->getGunslinger()->getGame()->getId()->toString(),
        ]);
    }

    /**
     * @param GameEvent $event
     *
     * @throws Exception
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     * @throws MessagePatternIsInvalidException
     * @throws \Psr\Cache\InvalidArgumentException
     */
    public function sendMessageWhenReadyToPlay(GameEvent $event): void
    {
        foreach ($this->messageBuilder->buildReadyToPlay() as $message) {
            $this->api->sendMessage(
                $event->getGame()->getChat()->getId(),
                $message,
                ParseMode::MARKDOWN
            );
            sleep(1);
        }

        $this->logger->info('event_subscriber.game.game.send_message_when_ready_to_play', [
            'game.id' => $event->getGame()->getId()->toString(),
        ]);
    }

    /**
     * @param GameEvent $event
     *
     * @throws Exception
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function sendMessageWhenPlayedWithNuclearBullet(GameEvent $event): void
    {
        if (!$event->getGame()->isPlayedWithNuclearBullet()) {
            return;
        }

        $this->api->sendMessage(
            $event->getGame()->getChat()->getId(),
            $this->messageBuilder->buildPlayedWithNuclearBullet(),
            ParseMode::MARKDOWN
        );

        $this->logger->info('event_subscriber.game.game.send_message_when_played_with_nuclear_bullet', [
            'game.id' => $event->getGame()->getId()->toString(),
        ]);
    }
}
