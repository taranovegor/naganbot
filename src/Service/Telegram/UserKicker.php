<?php
/**
 * This file is part of the naganbot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Service\Telegram;

use App\Constant\Common\Logger\Channel;
use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Service\Common\LoggerFactory;
use Psr\Log\LoggerInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Exception;

/**
 * Class UserKicker
 */
class UserKicker
{
    private BotApi $api;

    private LoggerInterface $logger;

    /**
     * UserKicker constructor.
     *
     * @param BotApi        $api
     * @param LoggerFactory $loggerFactory
     */
    public function __construct(BotApi $api, LoggerFactory $loggerFactory)
    {
        $this->api = $api;
        $this->logger = $loggerFactory->create(Channel::USER);
    }

    /**
     * @param Chat $chat
     * @param User $user
     * @param bool $unban
     *
     * @return bool
     */
    public function kick(Chat $chat, User $user, bool $unban = true): bool
    {
        try {
            $kicked = $this->api->kickChatMember(
                $chat->getId(),
                $user->getId(),
                $unban ? (new \DateTime('+1 minute'))->getTimestamp() : null
            );

            $this->logger->info('service.telegram.user_kicker.kick.kicked', [
                'chat.id' => $chat->getId(),
                'user.id' => $user->getId(),
            ]);

            return $kicked;
        } catch (Exception $e) {
            $this->logger->info('service.telegram.user_kicker.kick.was_not_kicked', [
                'chat.id' => $chat->getId(),
                'user.id' => $user->getId(),
                'exception.code' => $e->getCode(),
                'exception.message' => $e->getMessage(),
            ]);

            return false;
        }
    }
}
