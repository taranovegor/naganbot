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

use App\Repository\Telegram\ChatRepository;
use App\ValueObject\Telegram\Chat\Type;
use Doctrine\ORM\Mapping as ORM;

/**
 * Class Chat
 *
 * @ORM\Table(name="telegram_chat")
 * @ORM\Entity(repositoryClass=ChatRepository::class)
 *
 * @link https://core.telegram.org/bots/api#chat
 */
class Chat
{
    /**
     * @ORM\Id()
     * @ORM\Column(name="id", type="bigint", unique=true, nullable=false)
     */
    private int $id;

    /**
     * @ORM\Embedded(class=Type::class)
     */
    private Type $type;

    /**
     * @ORM\Column(name="title", type="string", nullable=true)
     */
    private ?string $title;

    /**
     * @ORM\Column(name="first_name", type="string", nullable=true)
     */
    private ?string $firstName;

    /**
     * @ORM\Column(name="last_name", type="string", nullable=true)
     */
    private ?string $lastName;

    /**
     * @ORM\Column(name="invite_link", type="string", nullable=true)
     */
    private ?string $inviteLink;

    /**
     * Chat constructor.
     *
     * @param int  $id
     * @param Type $type
     */
    public function __construct(int $id, Type $type)
    {
        $this->id = $id;
        $this->type = $type;
        $this->inviteLink = null;
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }


    /**
     * @return string
     */
    public function getType(): string
    {
        return $this->type;
    }

    /**
     * @param Type $type
     *
     * @return Chat
     */
    public function setType(Type $type): Chat
    {
        $this->type = $type;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getTitle(): ?string
    {
        return $this->title;
    }

    /**
     * @param null|string $title
     *
     * @return Chat
     */
    public function setTitle(?string $title): Chat
    {
        $this->title = $title;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getFirstName(): ?string
    {
        return $this->firstName;
    }

    /**
     * @param null|string $firstName
     *
     * @return Chat
     */
    public function setFirstName(?string $firstName): Chat
    {
        $this->firstName = $firstName;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getLastName(): ?string
    {
        return $this->lastName;
    }

    /**
     * @param null|string $lastName
     *
     * @return Chat
     */
    public function setLastName(?string $lastName): Chat
    {
        $this->lastName = $lastName;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getInviteLink(): ?string
    {
        return $this->inviteLink;
    }

    /**
     * @param null|string $inviteLink
     *
     * @return Chat
     */
    public function setInviteLink(?string $inviteLink): Chat
    {
        $this->inviteLink = $inviteLink;

        return $this;
    }
}
