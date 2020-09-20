<?php
/**
 * Copyright (C) 20.09.20 Egor Taranov
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

namespace App\Model\Translation\Intl;

/**
 * Class IntlMessageVariation
 */
class IntlMessageVariation
{
    private string $selector;

    /**
     * @var string|IntlVariableMessage
     */
    private $content;

    /**
     * IntlMessageVariation constructor.
     *
     * @param string $selector
     * @param        $content
     */
    public function __construct(string $selector, $content)
    {
        $this->selector = $selector;
        $this->content = $content;
    }

    /**
     * @return int|float|string
     */
    public function getSelector()
    {
        return $this->selector;
    }

    /**
     * @return IntlVariableMessage|string
     */
    public function getContent()
    {
        return $this->content;
    }

    /**
     * @return bool
     */
    public function isVariable(): bool
    {
        return $this->content instanceof IntlVariableMessage;
    }
}
