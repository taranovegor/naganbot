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

namespace App\ValueObject\Telegram\User;

use App\Entity\Telegram\User;

/**
 * Class UserGameStatistic
 */
class UserChatStatistic
{
    private User $user;

    private int $numberOfWins;

    /**
     * UserChatStatistic constructor.
     *
     * @param User $user
     * @param int  $numberOfWins
     */
    public function __construct(User $user, int $numberOfWins)
    {
        $this->user = $user;
        $this->numberOfWins = $numberOfWins;
    }

    /**
     * @return User
     */
    public function getUser(): User
    {
        return $this->user;
    }

    /**
     * @return int
     */
    public function getNumberOfWins(): int
    {
        return $this->numberOfWins;
    }
}
