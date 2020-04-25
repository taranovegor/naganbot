<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Telegram\Command;

use App\Exception\Game\NotEnoughGunslingersException;
use App\Exception\Game\ShotDeadNotFoundException;
use App\MessageBuilder\GameMessageBuilder;
use App\Manager\GameManager;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class JoinCommand
 */
class JoinCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrjoin';

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
        return $this->translator->trans('command.join');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @return void
     *
     * @throws ShotDeadNotFoundException
     * @throws \App\Exception\EntityNotFoundException
     * @throws \App\Exception\Game\AlreadyRegisteredInGameException
     * @throws \App\Exception\Game\GameIsAlreadyPlayedException
     * @throws \App\Exception\Game\NotFoundActiveGameException
     * @throws \Doctrine\ORM\ORMException
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     * @throws \Throwable
     */
    public function execute(BotApi $api, Update $update)
    {
        $this->gameManager->join(
            $update->getMessage()->getChat(),
            $update->getMessage()->getFrom()
        );

        try {
            $this->gameManager->play(
                $update->getMessage()->getChat()
            );
        } catch (NotEnoughGunslingersException $e) {
            // ignore
        }
    }
}
