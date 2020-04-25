<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventSubscriber;

use App\Event\GameEvent;
use App\Event\GunslingerEvent;
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
            GunslingerEvent::JOINED => [
                ['onGunslingerJoined', 0],
            ],
            GunslingerEvent::SHOT_HIMSELF => [
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
            $event->getGame()->getChat()->getId(),
            $this->gameMessageBuilder->buildCreate()->toString(),
            ParseMode::DEFAULT
        );
    }

    /**
     * @param GunslingerEvent $event
     *
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     */
    public function onGunslingerJoined(GunslingerEvent $event)
    {
        $this->api->sendMessage(
            $event->getGunslinger()->getGame()->getChat()->getId(),
            $this->gameMessageBuilder->buildJoin()->toString(),
            ParseMode::DEFAULT
        );
    }

    /**
     * @param GunslingerEvent $event
     *
     * @throws GunslingerWasNotKickedException
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     */
    public function onGunslingerShotHimself(GunslingerEvent $event)
    {
        $game = $event->getGunslinger()->getGame();

        /** @var Message $message */
        foreach ($this->gameMessageBuilder->buildPlay() as $message) {
            $this->api->sendMessage(
                $game->getChat()->getId(),
                $message->toString(),
                ParseMode::DEFAULT
            );
            sleep(1);
        }

        $this->api->sendMessage(
            $game->getChat()->getId(),
            $this->gameMessageBuilder->buildShotHimself(
                $event->getGunslinger()
            )->toString(),
            ParseMode::DEFAULT
        );

        try {
            $this->api->kickChatMember(
                $game->getChat()->getId(),
                $event->getGunslinger()->getUser()->getId()
            );
            $this->api->unbanChatMember(
                $game->getChat()->getId(),
                $event->getGunslinger()->getUser()->getId()
            );
        } catch (\TelegramBot\Api\Exception $e) {
            throw new GunslingerWasNotKickedException();
        }
    }
}
