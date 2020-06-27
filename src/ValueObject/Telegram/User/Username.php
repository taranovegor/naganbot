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

use App\ValueObject\NullableInterface;
use App\ValueObject\StringableInterface;
use Doctrine\ORM\Mapping as ORM;

/**
 * Class Username
 *
 * @ORM\Embeddable
 */
class Username implements StringableInterface, NullableInterface
{
    /**
     * @ORM\Column(name="username", type="string", nullable=true)
     */
    private ?string $username;

    /**
     * Username constructor.
     *
     * @param string|null $username
     */
    public function __construct(?string $username)
    {
        $this->username = $username;
    }

    /**
     * @return bool
     */
    public function isNull(): bool
    {
        return null === $this->username;
    }

    /**
     * @return string
     */
    public function toString(): string
    {
        return (string) $this->username;
    }

    /**
     * @return string
     */
    public function getMention(): string
    {
        return sprintf('@%s', $this->toString());
    }

    /**
     * @return string
     */
    public function __toString(): string
    {
        return $this->toString();
    }
}
