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

namespace App\Twig\Common;

use Twig\Extension\AbstractExtension;
use Twig\TwigFilter;

/**
 * Class StringExtension
 */
class StringExtension extends AbstractExtension
{
    /**
     * @inheritDoc
     */
    public function getFilters()
    {
        return [
            new TwigFilter('inline', [$this, 'inline']),
            new TwigFilter('line_break', [$this, 'lineBreak']),
        ];
    }

    /**
     * @param string $string
     *
     * @return string
     */
    public function inline(string $string): string
    {
        $string = preg_replace('!\s+!', ' ', $string);
        $string = str_replace(PHP_EOL, '', $string);
        $string = trim($string);

        return $string;
    }

    /**
     * @param string $string
     *
     * @return string
     */
    public function lineBreak(string $string): string
    {
        return $string.PHP_EOL;
    }
}
