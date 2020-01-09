<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Command;

use App\Service\GameManager;
use App\MessageBuilder\GameMessageBuilder;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
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
     * @throws \App\Exception\Game\GameIsAlreadyCreatedException
     * @throws \App\Exception\Game\GameIsAlreadyPlayedException
     * @throws \Doctrine\ORM\ORMException
     * @throws \App\Exception\EntityNotFoundException
     */
    public function execute(BotApi $api, Update $update)
    {
        $this->gameManager->create(
            $update->getMessage()->getChat(),
            $update->getMessage()->getFrom()
        );
    }
}
