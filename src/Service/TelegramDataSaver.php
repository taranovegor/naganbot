<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Service;

use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Exception\CatchableException;
use App\Exception\ChatOnlyException;
use App\Exception\EntityNotFoundException;
use App\Repository\Telegram\ChatRepository;
use App\Repository\Telegram\UserRepository;
use TelegramBot\Api\Types\Chat as TelegramChat;
use TelegramBot\Api\Types\MessageEntity;
use TelegramBot\Api\Types\Update;
use TelegramBot\Api\Types\User as TelegramUser;

/**
 * Class TelegramDataSaver
 */
class TelegramDataSaver
{
    /**
     * @var UserRepository
     */
    private UserRepository $userRepository;

    /**
     * @var ChatRepository
     */
    private ChatRepository $chatRepository;

    /**
     * @var Flusher
     */
    private Flusher $flusher;

    /**
     * TelegramDataSaver constructor.
     *
     * @param UserRepository $userRepository
     * @param ChatRepository $chatRepository
     * @param Flusher        $flusher
     */
    public function __construct(UserRepository $userRepository, ChatRepository $chatRepository, Flusher $flusher)
    {
        $this->userRepository = $userRepository;
        $this->chatRepository = $chatRepository;
        $this->flusher = $flusher;
    }

    /**
     * @param Update $update
     *
     * @throws CatchableException
     * @throws ChatOnlyException
     */
    public function handle(Update $update)
    {
        $message = $update->getMessage();
        if (null === $message) {
            return;
        }

        /** @var MessageEntity[] $entities */
        $entities = (array) $message->getEntities();
        if (0 === count($entities)) {
            return;
        }

        foreach ($entities as $entity) {
            if (MessageEntity::TYPE_BOT_COMMAND !== $entity->getType()) {
                continue;
            }

            if (false === strpos($message->getText(), '/rr')) {
                return;
            }

            if (\App\Model\Telegram\Type\Chat::TYPE_PRIVATE === $message->getChat()->getType()) {
                throw new ChatOnlyException();
            }

            try {
                $this->handleChat($message->getChat());
                $this->handleUser($message->getFrom());

                $this->flusher->flush();
            } catch (\Throwable $e) {
                throw new CatchableException(
                    'telegram_data_saver.handle.throwable',
                    false,
                    0,
                    $e
                );
            }

            break;
        }
    }

    /**
     * @param TelegramChat $chat
     *
     * @return Chat
     *
     * @throws \Doctrine\ORM\ORMException
     */
    protected function handleChat(TelegramChat $chat): Chat
    {
        try {
            $entity = $this->chatRepository->get(
                $chat->getId()
            )->updateFromChatType($chat);
        } catch (EntityNotFoundException $e) {
            $entity = Chat::createFromChatType($chat);
            $this->chatRepository->add($entity);
        }

        return $entity;
    }

    /**
     * @param TelegramUser $user
     *
     * @return User
     *
     * @throws \Doctrine\ORM\ORMException
     */
    protected function handleUser(TelegramUser $user): User
    {
        try {
            $entity = $this->userRepository->get(
                $user->getId()
            )->updateFromUserType($user);
        } catch (EntityNotFoundException $e) {
            $entity = User::createFromUserType($user);
            $this->userRepository->add($entity);
        }

        return $entity;
    }
}
