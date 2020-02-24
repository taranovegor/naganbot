<?php
/**
 * This file is part of the nagan-bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Telegram\Command;

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
 * Class TopCommand
 */
class TopCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrtop';

    /**
     * @var TranslatorInterface
     */
    private TranslatorInterface $translator;

    /**
     * @var ChatRepository
     */
    private ChatRepository $chatRepository;

    /**
     * @var UserRepository
     */
    private UserRepository $userRepository;

    /**
     * @var UserMessageBuilder
     */
    private UserMessageBuilder $topMessageBuilder;

    /**
     * StatCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param ChatRepository      $chatRepository
     * @param UserRepository      $userRepository
     * @param UserMessageBuilder  $topMessageBuilder
     */
    public function __construct(TranslatorInterface $translator, ChatRepository $chatRepository, UserRepository $userRepository, UserMessageBuilder $topMessageBuilder)
    {
        $this->translator = $translator;
        $this->chatRepository = $chatRepository;
        $this->userRepository = $userRepository;
        $this->topMessageBuilder = $topMessageBuilder;
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
        return $this->translator->trans('command.top');
    }

    /**
     * @inheritDoc
     *
     * @param BotApi $api
     * @param Update $update
     *
     * @throws \App\Exception\EntityNotFoundException
     * @throws \TelegramBot\Api\Exception
     * @throws \TelegramBot\Api\InvalidArgumentException
     */
    public function execute(BotApi $api, Update $update)
    {
        $chat = $this->chatRepository->get($update->getMessage()->getChat()->getId());
        $users = $this->userRepository->topByChat($chat);
        $api->sendMessage(
            $chat->getId(),
            $this->topMessageBuilder->buildTopUsers($users)->toString(),
            ParseMode::DEFAULT
        );
    }
}
