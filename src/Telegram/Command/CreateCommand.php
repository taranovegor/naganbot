<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Telegram\Command;

use App\Exception\EntityNotFoundException;
use App\Exception\Game\GameIsAlreadyCreatedException;
use App\Exception\Game\GameIsAlreadyPlayedException;
use App\Manager\GameManager;
use App\MessageBuilder\GameMessageBuilder;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Doctrine\ORM\ORMException;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class StartCommand
 */
class CreateCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrcreate';

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
        return $this->translator->trans('command.create');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @return void
     *
     * @throws GameIsAlreadyCreatedException
     * @throws GameIsAlreadyPlayedException
     * @throws ORMException
     * @throws EntityNotFoundException
     */
    public function execute(BotApi $api, Update $update)
    {
        $this->gameManager->create(
            $update->getMessage()->getChat(),
            $update->getMessage()->getFrom()
        );
    }
}
