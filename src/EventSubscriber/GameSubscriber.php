<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventSubscriber;

use App\Event\GameEvent;
use App\Event\GunslingerEvent;
use App\Exception\Game\FailedToScrollDrumException;
use App\Exception\Game\GameIsAlreadyPlayedException;
use App\Exception\Game\GunslingerWasNotKickedException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Game\ShotDeadNotFoundException;
use App\Manager\GameManager;
use App\MessageBuilder\GameMessageBuilder;
use App\Model\Message;
use App\Model\Telegram\ParseMode;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Exception;
use TelegramBot\Api\InvalidArgumentException;

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

    private GameManager $gameManager;

    /**
     * GameSubscriber constructor.
     *
     * @param GameMessageBuilder $gameMessageBuilder
     * @param BotApi             $api
     * @param GameManager        $gameManager
     */
    public function __construct(GameMessageBuilder $gameMessageBuilder, BotApi $api, GameManager $gameManager)
    {
        $this->gameMessageBuilder = $gameMessageBuilder;
        $this->api = $api;
        $this->gameManager = $gameManager;
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
                ['sendMessageWhenGunslingerJoined', 0],
                ['checkPossibilityOfStartingGameWhenGunslingerJoined', 0],
            ],
            GunslingerEvent::SHOT_HIMSELF => [
                ['onGunslingerShotHimself', 0],
            ],
        ];
    }

    /**
     * @param GameEvent $event
     *
     * @throws Exception
     * @throws InvalidArgumentException
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
     * @throws Exception
     * @throws InvalidArgumentException
     */
    public function sendMessageWhenGunslingerJoined(GunslingerEvent $event)
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
     * @throws FailedToScrollDrumException
     * @throws GameIsAlreadyPlayedException
     * @throws NotEnoughGunslingersException
     * @throws ShotDeadNotFoundException
     */
    public function checkPossibilityOfStartingGameWhenGunslingerJoined(GunslingerEvent $event): void
    {
        if ($this->gameManager->isEnoughPlayers($event->getGunslinger()->getGame())) {
            $this->gameManager->playGame($event->getGunslinger()->getGame());
        }
    }

    /**
     * @param GunslingerEvent $event
     *
     * @throws GunslingerWasNotKickedException
     * @throws Exception
     * @throws InvalidArgumentException
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
        } catch (Exception $e) {
            throw new GunslingerWasNotKickedException();
        }
    }
}
