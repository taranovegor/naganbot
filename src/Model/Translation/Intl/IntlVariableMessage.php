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
 * Class IntlVariableMessage
 */
class IntlVariableMessage implements \Countable, \Iterator
{
    public const TYPE_PLURAL = 'plural';
    public const TYPE_SELECT = 'select';

    private string $param;

    private string $type;

    /**
     * @var IntlMessageVariation[]
     */
    private array $variations;

    /**
     * IntlMessageFormatInformation constructor.
     *
     * @param string                 $param
     * @param string                 $type
     * @param IntlMessageVariation[] $variations
     */
    public function __construct(string $param, string $type, array $variations)
    {
        $this->param = $param;
        $this->type = $type;
        $this->variations = [];
        foreach ($variations as $variation) {
            $this->variations[$variation->getSelector()] = $variation;
        }
    }

    /**
     * @return string
     */
    public function getParam(): string
    {
        return $this->param;
    }

    /**
     * @return string
     */
    public function getType(): string
    {
        return $this->type;
    }

    /**
     * @return bool
     */
    public function isPlural(): bool
    {
        return self::TYPE_PLURAL === $this->type;
    }

    /**
     * @return bool
     */
    public function isSelect(): bool
    {
        return self::TYPE_SELECT === $this->type;
    }

    /**
     * @return IntlMessageVariation[]|array
     */
    public function getVariations(): array
    {
        return $this->variations;
    }

    /**
     * @param int|string $variation
     *
     * @return IntlMessageVariation
     */
    public function getVariation($variation): IntlMessageVariation
    {
        return $this->variations[$variation] ?? $this->variations['other'];
    }

    /**
     * @return IntlMessageVariation
     */
    public function first(): IntlMessageVariation
    {
        return $this->getVariation(array_key_first($this->getVariations()));
    }

    /**
     * @inheritDoc
     */
    public function count()
    {
        return count($this->variations);
    }

    /**
     * @inheritDoc
     */
    public function current()
    {
        return current($this->variations);
    }

    /**
     * @inheritDoc
     */
    public function next()
    {
        return next($this->variations);
    }

    /**
     * @inheritDoc
     */
    public function key()
    {
        return key($this->variations);
    }

    /**
     * @inheritDoc
     */
    public function valid()
    {
        return isset($this->variations[$this->key()]);
    }

    /**
     * @inheritDoc
     */
    public function rewind()
    {
        reset($this->variations);
    }
}
