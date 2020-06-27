<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Telegram\Command;

use App\Exception\EntityNotFoundException;
use App\Exception\Game\FailedToScrollDrumException;
use App\Exception\Game\GameIsAlreadyPlayedException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Game\ShotDeadNotFoundException;
use App\Manager\ChatManager;
use App\Manager\GameManager;
use App\MessageBuilder\GameMessageBuilder;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class ForceCommand
 */
class ForceCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrforce';

    /**
     * @var TranslatorInterface
     */
    private TranslatorInterface $translator;

    /**
     * @var GameManager
     */
    private GameManager $gameManager;

    /**
     * @var GameMessageBuilder
     */
    private GameMessageBuilder $gameMessageBuilder;

    /**
     * @var ChatManager
     */
    private ChatManager $chatManager;

    /**
     * StartCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param GameManager         $gameManager
     * @param GameMessageBuilder  $gameMessageBuilder
     * @param ChatManager         $chatManager
     */
    public function __construct(TranslatorInterface $translator, GameManager $gameManager, GameMessageBuilder $gameMessageBuilder, ChatManager $chatManager)
    {
        $this->translator = $translator;
        $this->gameManager = $gameManager;
        $this->gameMessageBuilder = $gameMessageBuilder;
        $this->chatManager = $chatManager;
    }

    /**
     * @inheritDoc
     */
    public function getName()
    {
        return self::COMMAND;
    }

    /**
     * @inheritDoc
     */
    public function getDescription()
    {
        return $this->translator->trans('command.force');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @return void
     *
     * @throws ShotDeadNotFoundException
     * @throws EntityNotFoundException
     * @throws FailedToScrollDrumException
     */
    public function execute(BotApi $api, Update $update)
    {
        try {
            $chat = $this->chatManager->get((int) $update->getMessage()->getChat()->getId());
            $game = $this->gameManager->getLatestByChat($chat);
            $this->gameManager->playGame($game);
        } catch (NotEnoughGunslingersException | GameIsAlreadyPlayedException $e) {
            // ignore
        }
    }
}
