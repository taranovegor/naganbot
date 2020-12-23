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
use App\Event\Game\GunslingerEvent;
use App\Exception\Game\GunslingerWasNotKickedException;
use App\Service\Common\LoggerFactory;
use App\Service\MessageBuilder\Game\GunslingerMessageBuilder;
use App\Service\Telegram\UserKicker;
use Psr\Log\LoggerInterface;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Exception;
use TelegramBot\Api\InvalidArgumentException;
use Twig\Error\LoaderError;
use Twig\Error\RuntimeError;
use Twig\Error\SyntaxError;

/**
 * Class GunslingerSubscriber
 */
class GunslingerSubscriber implements EventSubscriberInterface
{
    private GunslingerMessageBuilder $messageBuilder;

    private BotApi $api;

    private LoggerInterface $logger;

    private UserKicker $kicker;

    /**
     * GunslingerSubscriber constructor.
     *
     * @param GunslingerMessageBuilder $messageBuilder
     * @param BotApi                   $api
     * @param LoggerFactory            $loggerFactory
     * @param UserKicker               $kicker
     */
    public function __construct(GunslingerMessageBuilder $messageBuilder, BotApi $api, LoggerFactory $loggerFactory, UserKicker $kicker)
    {
        $this->messageBuilder = $messageBuilder;
        $this->api = $api;
        $this->logger = $loggerFactory->create(Channel::GUNSLINGER);
        $this->kicker = $kicker;
    }

    /**
     * @inheritDoc
     */
    public static function getSubscribedEvents()
    {
        return [
            GunslingerEvent::JOINED_TO_GAME => [
                ['sendMessageWhenJoinedToGame', 0],
            ],
            GunslingerEvent::SHOT_HIMSELF => [
                ['sendMessageWhenShotHimself', 0],
            ],
        ];
    }

    /**
     * @param GunslingerEvent $event
     *
     * @throws Exception
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function sendMessageWhenJoinedToGame(GunslingerEvent $event): void
    {
        $gunslinger = $event->getGunslinger();

        if ($gunslinger->getGame()->isPlayed()) {
            $this->logger->info('event_subscriber.game.gunslinger.send_message_when_joined_to_game.ignored_game_already_played', [
                'gunslinger.id' => $event->getGunslinger()->getId()->toString(),
            ]);

            return;
        }

        if ($gunslinger->getUser()->isEquals($gunslinger->getGame()->getOwner())) {
            $this->logger->info('event_subscriber.game.gunslinger.send_message_when_joined_to_game.ignored_gunslinger_is_owner', [
                'gunslinger.id' => $event->getGunslinger()->getId()->toString(),
            ]);

            return;
        }

        $this->api->sendMessage(
            $event->getGunslinger()->getGame()->getChat()->getId(),
            $this->messageBuilder->buildJoin(),
            ParseMode::MARKDOWN
        );

        $this->logger->info('event_subscriber.game.gunslinger.send_message_when_joined_to_game', [
            'gunslinger.id' => $event->getGunslinger()->getId()->toString(),
        ]);
    }

    /**
     * @param GunslingerEvent $event
     *
     * @throws Exception
     * @throws GunslingerWasNotKickedException
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function sendMessageWhenShotHimself(GunslingerEvent $event): void
    {
        $gunslinger = $event->getGunslinger();

        $this->api->sendMessage(
            $gunslinger->getGame()->getChat()->getId(),
            $this->messageBuilder->buildShotHimself($gunslinger),
            ParseMode::MARKDOWN
        );

        $kicked = $this->kicker->kick(
            $gunslinger->getGame()->getChat(),
            $gunslinger->getUser()
        );

        try {
            if ($kicked) {
                $this->logger->info('event_subscriber.game.gunslinger.send_message_when_shot_himself.was_kicked', [
                    'gunslinger.id' => $event->getGunslinger()->getId()->toString(),
                    'user.id' => $event->getGunslinger()->getUser()->getId(),
                ]);
            } else {
                $this->logger->info('event_subscriber.game.gunslinger.send_message_when_shot_himself.was_not_kicked', [
                    'gunslinger.id' => $event->getGunslinger()->getId()->toString(),
                    'user.id' => $event->getGunslinger()->getUser()->getId(),
                ]);

                throw new GunslingerWasNotKickedException();
            }
        } finally {
            $this->logger->info('event_subscriber.game.gunslinger.send_message_when_shot_himself', [
                'gunslinger.id' => $event->getGunslinger()->getId()->toString(),
            ]);
        }
    }
}
