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

namespace App\Event\Game;

use App\Entity\Game\Gunslinger;
use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class GunslingerEvent
 */
class GunslingerEvent extends Event
{
    public const JOINED_TO_GAME = 'gunslinger.joined_to_game';
    public const SHOT_HIMSELF = 'gunslinger.shot_himself';

    /**
     * @var Gunslinger
     */
    private Gunslinger $gunslinger;

    /**
     * GunslingerEvent constructor.
     *
     * @param Gunslinger $gunslinger
     */
    public function __construct(Gunslinger $gunslinger)
    {
        $this->gunslinger = $gunslinger;
    }

    /**
     * @return Gunslinger
     */
    public function getGunslinger(): Gunslinger
    {
        return $this->gunslinger;
    }
}
