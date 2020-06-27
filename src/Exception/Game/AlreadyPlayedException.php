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

namespace App\Exception\Game;

use App\Exception\ParameterizedExceptionMessageInterface;

/**
 * Class AlreadyPlayedException
 */
class AlreadyPlayedException extends \Exception implements ParameterizedExceptionMessageInterface
{
    protected $message = 'game.already_played';

    /**
     * @return array
     */
    public function getMessageParameters(): array
    {
        return [
            'remaining_time' => (new \DateTime())->diff((new \DateTime('+1 day'))->setTime(0, 0, 0)),
        ];
    }
}
