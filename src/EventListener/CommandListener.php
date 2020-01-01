<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventListener;

use App\Entity\Telegram\Chat;
use App\Model\Telegram\Type\Chat as ChatModel;
use App\Entity\Telegram\User;
use App\Exception\CatchableException;
use App\Exception\ChatOnlyException;
use App\Exception\EntityNotFoundException;
use App\Repository\Telegram\ChatRepository;
use App\Repository\Telegram\UserRepository;
use App\Service\Flusher;
use BoShurik\TelegramBotBundle\Event\Telegram\UpdateEvent;
use Psr\Log\LoggerInterface;
use TelegramBot\Api\Types\MessageEntity;

/**
 * Class CommandSubscriber
 */
class CommandListener
{
    /**
     * @var ChatRepository
     */
    private ChatRepository $chatRepository;

    /**
     * @var UserRepository
     */
    private UserRepository $userRepository;

    /**
     * @var Flusher
     */
    private Flusher $flusher;

    /**
     * @var LoggerInterface
     */
    private LoggerInterface $logger;

    /**
     * CommandSubscriber constructor.
     *
     * @param ChatRepository  $chatRepository
     * @param UserRepository  $userRepository
     * @param Flusher         $flusher
     * @param LoggerInterface $logger
     */
    public function __construct(ChatRepository $chatRepository, UserRepository $userRepository, Flusher $flusher, LoggerInterface $logger)
    {
        $this->chatRepository = $chatRepository;
        $this->userRepository = $userRepository;
        $this->flusher = $flusher;
        $this->logger = $logger;
    }

    /**
     * @param UpdateEvent $event
     *
     * @throws CatchableException
     * @throws ChatOnlyException
     */
    public function onUpdate(UpdateEvent $event)
    {
        $message = $event->getUpdate()->getMessage();
        if (null === $message) {
            return;
        }

        /** @var MessageEntity[] $entities */
        $entities = $message->getEntities();
        if (null === $entities) {
            return;
        }

        foreach ($entities as $entity) {
            if (MessageEntity::TYPE_BOT_COMMAND !== $entity->getType()) {
                continue;
            }

            if (false === strpos($message->getText(), '/rr')) {
                continue;
            }

            if (ChatModel::TYPE_PRIVATE === $message->getChat()->getType()) {
                throw new ChatOnlyException();
            }

            try {
                try {
                    $this->chatRepository->get(
                        $message->getChat()->getId()
                    )->updateFromChatType($message->getChat());
                } catch (EntityNotFoundException $e) {
                    $this->chatRepository->add(Chat::createFromChatType($message->getChat()));
                }

                try {
                    $this->userRepository->get(
                        $message->getFrom()->getId()
                    )->updateFromUserType($message->getFrom());
                } catch (EntityNotFoundException $e) {
                    $this->userRepository->add(User::createFromUserType($message->getFrom()));
                }

                $this->flusher->flush();
            } catch (\Throwable $e) {
                throw new CatchableException(
                    'command_subscriber.check_update_for_any_command.throwable',
                    false,
                    0,
                    $e
                );
            }

            return;
        }
    }
}
