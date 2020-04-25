<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Telegram\Command;

use App\Exception\Game\GameIsAlreadyPlayedException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Game\ShotDeadNotFoundException;
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
     * StartCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param GameManager         $gameManager
     * @param GameMessageBuilder  $gameMessageBuilder
     */
    public function __construct(TranslatorInterface $translator, GameManager $gameManager, GameMessageBuilder $gameMessageBuilder)
    {
        $this->translator = $translator;
        $this->gameManager = $gameManager;
        $this->gameMessageBuilder = $gameMessageBuilder;
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
     * @throws \App\Exception\EntityNotFoundException
     * @throws \App\Exception\Game\FailedToScrollDrumException
     */
    public function execute(BotApi $api, Update $update)
    {
        try {
            $this->gameManager->play(
                $update->getMessage()->getChat()
            );
        } catch (NotEnoughGunslingersException | GameIsAlreadyPlayedException $e) {
            // ignore
        }
    }
}
