<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Command;

use App\Entity\Game;
use App\Exception\EntityNotFoundException;
use App\MessageBuilder\GameMessageBuilder;
use App\Model\Telegram\ParseMode;
use App\Repository\ShotHimselfRepository;
use App\Service\GameManager;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class JoinedCommand
 */
class JoinedCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrjoined';

    /**
     * @var TranslatorInterface
     */
    private TranslatorInterface $translator;

    /**
     * @var GameManager
     */
    private GameManager $gameManager;

    /**
     * @var ShotHimselfRepository
     */
    private ShotHimselfRepository $shotHimselfRepository;

    /**
     * @var GameMessageBuilder
     */
    private GameMessageBuilder $gameMessageBuilder;

    /**
     * StartCommand constructor.
     *
     * @param TranslatorInterface   $translator
     * @param GameManager           $gameManager
     * @param ShotHimselfRepository $shotHimselfRepository
     * @param GameMessageBuilder    $gameMessageBuilder
     */
    public function __construct(TranslatorInterface $translator, GameManager $gameManager, ShotHimselfRepository $shotHimselfRepository, GameMessageBuilder $gameMessageBuilder)
    {
        $this->translator = $translator;
        $this->gameManager = $gameManager;
        $this->shotHimselfRepository = $shotHimselfRepository;
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
        return $this->translator->trans('command.joined');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @return void
     *
     * @throws \TelegramBot\Api\InvalidArgumentException
     * @throws \TelegramBot\Api\Exception
     */
    public function execute(BotApi $api, Update $update)
    {
        try {
            $gunslingers = $this->gameManager->joined(
                $update->getMessage()->getChat()
            );
        } catch (EntityNotFoundException $e) {
            return;
        }

        /** @var Game $game */
        $game = $gunslingers[0]->getGame();

        $message = $this->gameMessageBuilder->buildJoined(
            $game,
            $this->shotHimselfRepository->findByGameTable($game)
        );

        $api->sendMessage(
            $game->getChat()->getId(),
            $message->toString(),
            ParseMode::DEFAULT
        );
    }
}
