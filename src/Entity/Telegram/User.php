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

namespace App\Entity\Telegram;

use App\Repository\Telegram\UserRepository;
use App\ValueObject\Telegram\User\LanguageCode;
use App\ValueObject\Telegram\User\Username;
use Doctrine\ORM\Mapping as ORM;

/**
 * Class User
 *
 * @ORM\Table(name="telegram_user")
 * @ORM\Entity(repositoryClass=UserRepository::class)
 *
 * @link https://core.telegram.org/bots/api#user
 */
class User
{
    /**
     * @ORM\Id()
     * @ORM\Column(name="id", type="integer")
     */
    private int $id;

    /**
     * @ORM\Column(name="is_bot", type="boolean")
     */
    private bool $isBot;

    /**
     * @ORM\Column(name="first_name", type="string")
     */
    private string $firstName;

    /**
     * @ORM\Column(name="last_name", type="string", nullable=true)
     */
    private ?string $lastName;

    /**
     * @ORM\Embedded(class=Username::class)
     */
    private Username $username;

    /**
     * @ORM\Embedded(class=LanguageCode::class)
     */
    private LanguageCode $languageCode;

    /**
     * User constructor.
     *
     * @param int    $id
     * @param string $firstName
     * @param bool   $isBot
     */
    public function __construct(int $id, string $firstName, bool $isBot)
    {
        $this->id = $id;
        $this->firstName = $firstName;
        $this->isBot = $isBot;
        $this->username = new Username(null);
        $this->languageCode = new LanguageCode(null);
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * @return bool
     */
    public function isBot(): bool
    {
        return $this->isBot;
    }

    /**
     * @param bool $isBot
     *
     * @return User
     */
    public function setIsBot(bool $isBot): User
    {
        $this->isBot = $isBot;

        return $this;
    }

    /**
     * @return string
     */
    public function getFirstName(): string
    {
        return $this->firstName;
    }

    /**
     * @param string $firstName
     *
     * @return User
     */
    public function setFirstName(string $firstName): User
    {
        $this->firstName = $firstName;

        return $this;
    }

    /**
     * @return string|null
     */
    public function getLastName(): ?string
    {
        return $this->lastName;
    }

    /**
     * @param string|null $lastName
     *
     * @return User
     */
    public function setLastName(?string $lastName): User
    {
        $this->lastName = $lastName;

        return $this;
    }

    /**
     * @return Username
     */
    public function getUsername(): Username
    {
        return $this->username;
    }

    /**
     * @param Username $username
     *
     * @return User
     */
    public function setUsername(Username $username): User
    {
        $this->username = $username;

        return $this;
    }

    /**
     * @return LanguageCode
     */
    public function getLanguageCode(): LanguageCode
    {
        return $this->languageCode;
    }

    /**
     * @param LanguageCode $languageCode
     *
     * @return User
     */
    public function setLanguageCode(LanguageCode $languageCode): User
    {
        $this->languageCode = $languageCode;

        return $this;
    }

    /**
     * @param User $user
     *
     * @return bool
     */
    public function isEquals(User $user): bool
    {
        return $this->getId() === $user->getId();
    }
}
