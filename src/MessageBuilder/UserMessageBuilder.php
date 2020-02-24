<?php
/**
 * This file is part of the nagan-bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\MessageBuilder;

use App\Model\Message;
use App\Model\Telegram\TopUser;
use App\Model\Username;

/**
 * Class ChatMessageBuilder
 */
class UserMessageBuilder extends AbstractTranslatableMessageBuilder
{
    /**
     * @param array|TopUser[] $users
     *
     * @return Message
     */
    public function buildTopUsers(array $users): Message
    {
        $this->message->addLine($this->translator->trans('user.top.title', [
            '%count%' => count($users),
        ], 'messages'));

        foreach ($users as $number => $user) {
            $this->message->addLine($this->translator->trans('user.top.item', [
                '%num%' => $number + 1,
                '%user%' => Username::fromUser($user->getUser())->toString(),
                '%count%' => $user->getNumberOfWins(),
            ], 'messages'));
        }

        return $this->message;
    }
}
