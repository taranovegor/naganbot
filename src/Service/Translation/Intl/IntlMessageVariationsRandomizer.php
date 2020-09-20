<?php
/**
 * Copyright (C) 26.09.20 Egor Taranov
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

namespace App\Service\Translation\Intl;

use App\Constant\Common\Translation\Domain;
use App\Exception\Translation\MessagePatternIsInvalidException;
use App\Model\Translation\Intl\IntlMessageVariation;
use App\Model\Translation\Intl\IntlVariableMessage;
use Psr\Cache\InvalidArgumentException;

/**
 * Class IntlMessageVariationsRandomizer
 */
class IntlMessageVariationsRandomizer
{
    private IntlMessageVariationsProvider $provider;

    /**
     * IntlMessageVariationsRandomizer constructor.
     *
     * @param IntlMessageVariationsProvider $provider
     */
    public function __construct(IntlMessageVariationsProvider $provider)
    {
        $this->provider = $provider;
    }

    /**
     * @param string      $id         translation key containing intl string
     * @param array       $properties an array of keys
     *                                all keys obtained with a null value will be randomized
     *                                all keys with non-null values will be returned as is
     * @param string      $domain
     * @param string|null $locale
     *
     * @return array
     *
     * @throws InvalidArgumentException
     * @throws MessagePatternIsInvalidException
     */
    public function rand(string $id, array $properties = [], string $domain = Domain::DEFAULT, string $locale = null): array
    {
        $message = $this->provider->provide($id, $domain, $locale);
        if (!$message instanceof IntlVariableMessage) {
            return $properties;
        }

        $this->randIntlVariableMessage($message, $properties);

        return $properties;
    }

    /**
     * @param IntlVariableMessage $variable
     * @param array               $properties
     *
     * @return IntlVariableMessage
     */
    public function randIntlVariableMessage(IntlVariableMessage $variable, array &$properties = []): IntlVariableMessage
    {
        if (false === array_key_exists($variable->getParam(), $properties)) {
            return $variable;
        }

        $selector = $properties[$variable->getParam()] ?? array_rand($variable->getVariations());
        $variation = $variable->getVariation($selector);
        $properties = array_merge($properties, [$variable->getParam() => $selector]);
        if ($variation->getContent() instanceof IntlVariableMessage) {
            $variation = new IntlMessageVariation(
                $variation->getSelector(),
                $this->randIntlVariableMessage(
                    $variation->getContent(),
                    $properties
                )
            );
        }

        return new IntlVariableMessage($variable->getParam(), $variable->getType(), [$variation]);
    }
}
