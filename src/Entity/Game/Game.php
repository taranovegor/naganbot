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

namespace App\Entity\Game;

use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Repository\Game\GameRepository;
use DateTime;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\ORM\Mapping as ORM;
use Ramsey\Uuid\Uuid;
use Ramsey\Uuid\UuidInterface;

/**
 * Class GameTable
 *
 * @ORM\Table(name="game_game")
 * @ORM\Entity(repositoryClass=GameRepository::class)
 */
class Game
{
    /**
     * @ORM\Id()
     * @ORM\Column(name="id", type="uuid")
     */
    private UuidInterface $id;

    /**
     * @ORM\ManyToOne(targetEntity=Chat::class, cascade={"persist"})
     * @ORM\JoinColumn(name="chat_id", referencedColumnName="id", nullable=false)
     */
    private Chat $chat;

    /**
     * @ORM\ManyToOne(targetEntity=User::class, cascade={"persist"})
     * @ORM\JoinColumn(name="owner_user_id", referencedColumnName="id", nullable=false)
     */
    private User $owner;

    /**
     * @var Gunslinger[]|Collection
     *
     * @ORM\OneToMany(targetEntity=Gunslinger::class, cascade={"persist"}, mappedBy="game")
     * @ORM\OrderBy({"joinedAt" = "ASC"})
     */
    private Collection $gunslingers;

    /**
     * @ORM\Column(name="played_with_nuclear_bullet", type="boolean")
     */
    private bool $playedWithNuclearBullet;

    /**
     * @ORM\Column(name="created_at", type="datetime")
     */
    private DateTime $createdAt;

    /**
     * @ORM\Column(name="played_at", type="datetime", nullable=true)
     */
    private ?DateTime $playedAt;

    /**
     * GameTable constructor.
     *
     * @param Chat $chat
     * @param User $owner
     */
    public function __construct(Chat $chat, User $owner)
    {
        $this->id = Uuid::uuid4();
        $this->chat = $chat;
        $this->owner = $owner;
        $this->gunslingers = new ArrayCollection();
        $this->playedWithNuclearBullet = false;
        $this->createdAt = new DateTime();
        $this->playedAt = null;
    }

    /**
     * @return UuidInterface
     */
    public function getId(): UuidInterface
    {
        return $this->id;
    }

    /**
     * @return Chat
     */
    public function getChat(): Chat
    {
        return $this->chat;
    }

    /**
     * @return User
     */
    public function getOwner(): User
    {
        return $this->owner;
    }

    /**
     * @return Gunslinger[]|Collection
     */
    public function getGunslingers(): Collection
    {
        return $this->gunslingers;
    }

    /**
     * @param Gunslinger $gunslinger
     *
     * @return Game
     */
    public function addGunslinger(Gunslinger $gunslinger): Game
    {
        if (!$this->gunslingers->contains($gunslinger)) {
            $this->gunslingers->add($gunslinger);
        }

        return $this;
    }

    /**
     * @return bool
     */
    public function isPlayedWithNuclearBullet(): bool
    {
        return $this->playedWithNuclearBullet;
    }

    /**
     * @return Game
     */
    public function markAsPlayedWithNuclearBullet(): Game
    {
        if (!$this->isPlayed()) {
            $this->playedWithNuclearBullet = true;
        }

        return $this;
    }

    /**
     * @return DateTime
     */
    public function getCreatedAt(): DateTime
    {
        return $this->createdAt;
    }

    /**
     * @return DateTime|null
     */
    public function getPlayedAt(): ?DateTime
    {
        return $this->playedAt;
    }

    /**
     * @return bool
     */
    public function isPlayed(): bool
    {
        return $this->playedAt instanceof DateTime;
    }

    /**
     * @return Game
     */
    public function markAsPlayed(): Game
    {
        if (null === $this->playedAt) {
            $this->playedAt = new DateTime();
        }

        return $this;
    }

    /**
     * @return bool
     */
    public function isCreatedToday(): bool
    {
        return 0 === (int) (clone $this->createdAt->modify('today'))
            ->diff(new DateTime())
            ->format('%R%a')
        ;
    }
}
