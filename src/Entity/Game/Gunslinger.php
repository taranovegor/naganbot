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

use App\Entity\Telegram\User;
use App\Repository\Game\GunslingerRepository;
use DateTime;
use Doctrine\ORM\Mapping as ORM;
use Ramsey\Uuid\Uuid;
use Ramsey\Uuid\UuidInterface;

/**
 * Class Gunslinger
 *
 * @ORM\Table(name="game_gunslinger")
 * @ORM\Entity(repositoryClass=GunslingerRepository::class)
 */
class Gunslinger
{
    /**
     * @var UuidInterface
     *
     * @ORM\Id()
     * @ORM\Column(name="id", type="uuid_binary")
     */
    private UuidInterface $id;

    /**
     * @var Game
     *
     * @ORM\ManyToOne(targetEntity=Game::class, cascade={"persist"}, inversedBy="gunslingers")
     * @ORM\JoinColumn(name="game_id", referencedColumnName="id", nullable=false)
     */
    private Game $game;

    /**
     * @var User
     *
     * @ORM\ManyToOne(targetEntity=User::class, cascade={"persist"})
     * @ORM\JoinColumn(name="user_id", referencedColumnName="id", nullable=false)
     */
    private User $user;

    /**
     * @var DateTime
     *
     * @ORM\Column(name="joined_at", type="datetime")
     */
    private DateTime $joinedAt;

    /**
     * @var bool
     *
     * @ORM\Column(name="shot_himself", type="boolean")
     */
    private bool $shotHimself;

    /**
     * Gunslinger constructor.
     *
     * @param Game $game
     * @param User $user
     */
    public function __construct(Game $game, User $user)
    {
        $this->id = Uuid::uuid4();
        $this->game = $game;
        $this->user = $user;
        $this->joinedAt = new DateTime();
        $this->shotHimself = false;
    }

    /**
     * @return UuidInterface
     */
    public function getId(): UuidInterface
    {
        return $this->id;
    }

    /**
     * @return Game
     */
    public function getGame(): Game
    {
        return $this->game;
    }

    /**
     * @return User
     */
    public function getUser(): User
    {
        return $this->user;
    }

    /**
     * @return DateTime
     */
    public function getJoinedAt(): DateTime
    {
        return $this->joinedAt;
    }

    /**
     * @return bool
     */
    public function isShotHimself(): bool
    {
        return $this->shotHimself;
    }

    /**
     * @return Gunslinger
     */
    public function shot(): Gunslinger
    {
        $this->shotHimself = true;

        return $this;
    }
}
