<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Telegram\Command;

use App\Exception\EntityNotFoundException;
use App\MessageBuilder\UserMessageBuilder;
use App\Model\Telegram\ParseMode;
use App\Repository\Telegram\ChatRepository;
use App\Repository\Telegram\UserRepository;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class MeCommand
 */
class MeCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrme';

    private TranslatorInterface $translator;

    private ChatRepository $chatRepository;

    private UserRepository $userRepository;

    private UserMessageBuilder $messageBuilder;

    /**
     * MeCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param ChatRepository      $chatRepository
     * @param UserRepository      $userRepository
     */
    public function __construct(TranslatorInterface $translator, ChatRepository $chatRepository, UserRepository $userRepository, UserMessageBuilder $messageBuilder)
    {
        $this->translator = $translator;
        $this->chatRepository = $chatRepository;
        $this->userRepository = $userRepository;
        $this->messageBuilder = $messageBuilder;
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
        return $this->translator->trans('command.me');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @return mixed|void
     *
     * @throws EntityNotFoundException
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     */
    public function execute(BotApi $api, Update $update)
    {
        $user = $this->userRepository->get($update->getMessage()->getFrom()->getId());
        $chat = $this->chatRepository->get($update->getMessage()->getChat()->getId());
        $statistic = $this->userRepository->getStatisticByUserInChat($user, $chat);
        $message = $this->messageBuilder->buildMe($statistic);

        $api->sendMessage(
            $chat->getId(),
            $message->toString(),
            ParseMode::DEFAULT
        );
    }
}
