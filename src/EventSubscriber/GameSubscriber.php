<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventSubscriber;

use App\Event\GameEvent;
use App\Event\GunslingerJoinedEvent;
use App\Event\GunslingerShotHimselfEvent;
use App\Exception\Game\GunslingerWasNotKickedException;
use App\MessageBuilder\GameMessageBuilder;
use App\Model\Message;
use App\Model\Telegram\ParseMode;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use TelegramBot\Api\BotApi;

/**
 * Class GameSubscriber
 */
class GameSubscriber implements EventSubscriberInterface
{
    /**
     * @var GameMessageBuilder
     */
    private GameMessageBuilder $gameMessageBuilder;

    /**
     * @var BotApi
     */
    private BotApi $api;

    /**
     * GameSubscriber constructor.
     *
     * @param GameMessageBuilder $gameMessageBuilder
     * @param BotApi             $api
     */
    public function __construct(GameMessageBuilder $gameMessageBuilder, BotApi $api)
    {
        $this->gameMessageBuilder = $gameMessageBuilder;
        $this->api = $api;
    }

    /**
     * @inheritDoc
     */
    public static function getSubscribedEvents()
    {
        return [
            GameEvent::CREATED => [
                ['onGameCreated', 0],
            ],
            GunslingerJoinedEvent::EVENT => [
                ['onGunslingerJoined', 0],
            ],
            GunslingerShotHimselfEvent::EVENT => [
                ['onGunslingerShotHimself', 0],
            ],
        ];
    }

    /**
     * @param GameEvent $event
     *
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     */
    public function onGameCreated(GameEvent $event)
    {
        $this->api->sendMessage(
            $event->getGameTable()->getChat()->getId(),
            $this->gameMessageBuilder->buildCreate()->toString(),
            ParseMode::DEFAULT
        );
    }

    /**
     * @param GunslingerJoinedEvent $event
     *
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     */
    public function onGunslingerJoined(GunslingerJoinedEvent $event)
    {
        $this->api->sendMessage(
            $event->getGameTable()->getChat()->getId(),
            $this->gameMessageBuilder->buildJoin()->toString(),
            ParseMode::DEFAULT
        );
    }

    /**
     * @param GunslingerShotHimselfEvent $event
     *
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     * @throws GunslingerWasNotKickedException
     */
    public function onGunslingerShotHimself(GunslingerShotHimselfEvent $event)
    {
        /** @var Message $message */
        foreach ($this->gameMessageBuilder->buildPlay() as $message) {
            $this->api->sendMessage(
                $event->getGameTable()->getChat()->getId(),
                $message->toString(),
                ParseMode::DEFAULT
            );
            sleep(1);
        }

        $this->api->sendMessage(
            $event->getGameTable()->getChat()->getId(),
            $this->gameMessageBuilder->buildShotHimself(
                $event->getShotHimself()
            )->toString(),
            ParseMode::DEFAULT
        );

        try {
            $this->api->kickChatMember(
                $event->getGameTable()->getChat()->getId(),
                $event->getShotHimself()->getGunslinger()->getUser()->getId()
            );
            $this->api->unbanChatMember(
                $event->getGameTable()->getChat()->getId(),
                $event->getShotHimself()->getGunslinger()->getUser()->getId()
            );
        } catch (\TelegramBot\Api\Exception $e) {
            throw new GunslingerWasNotKickedException();
        }
    }
}
