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

namespace App\Twig\Telegram;

use App\Entity\Telegram\User;
use Twig\Extension\AbstractExtension;
use Twig\TwigFilter;

/**
 * Class UserExtension
 */
class UserExtension extends AbstractExtension
{
    /**
     * @inheritDoc
     */
    public function getFilters()
    {
        return [
            new TwigFilter('fullName', [$this, 'fullName']),
            new TwigFilter('username', [$this, 'username']),
        ];
    }

    /**
     * @param User $user
     *
     * @return string
     */
    public function fullName(User $user): string
    {
        return trim(sprintf('%s %s', $user->getFirstName(), $user->getLastName()));
    }

    /**
     * @param User $user
     *
     * @return string
     */
    public function username(User $user): string
    {
        return $user->getUsername()->toString() ?? $this->fullName($user);
    }
}
