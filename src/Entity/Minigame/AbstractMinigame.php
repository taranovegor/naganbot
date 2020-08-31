<?php
/**
 * Copyright (C) 30.08.20 Egor Taranov
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

namespace App\Entity\Minigame;

use Doctrine\ORM\Mapping as ORM;
use Doctrine\ORM\Mapping\Entity;
use Ramsey\Uuid\UuidInterface;

/**
 * Class AbstractMinigame
 *
 * @Entity
 * @ORM\InheritanceType("JOINED")
 * @ORM\DiscriminatorColumn(name="type")
 * @ORM\DiscriminatorMap({
 *     AbstractMinigame::TYPE_NAGAN_OR_MAUSER = NaganOrMauser::class
 * })
 */
abstract class AbstractMinigame
{
    public const TYPE_NAGAN_OR_MAUSER = 'nagan_or_mauser';

    /**
     * @ORM\Id()
     * @ORM\Column(name="id", type="uuid")
     */
    private ?UuidInterface $id;

    /**
     * @ORM\Column(name="type", type="string", length=16)
     */
    private string $type;
}
