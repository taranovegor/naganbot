<?php
/**
 * Copyright (C) 14.08.20 Egor Taranov
 * This file is part of Nagan bot <https://github.com/taranovegor/nagan-bot>.
 *
 * Nagan bot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Nagan bot is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Nagan bot.  If not, see <http://www.gnu.org/licenses/>.
 */

namespace App\Event\Telegram;

use App\Entity\Telegram\Chat;
use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class ChatEvent
 */
class ChatEvent extends Event
{
    public const CREATED = 'telegram.chat.created';
    public const UPDATED = 'telegram.chat.updated';

    private Chat $chat;

    /**
     * ChatEvent constructor.
     *
     * @param Chat $chat
     */
    public function __construct(Chat $chat)
    {
        $this->chat = $chat;
    }

    /**
     * @return Chat
     */
    public function getChat(): Chat
    {
        return $this->chat;
    }
}
