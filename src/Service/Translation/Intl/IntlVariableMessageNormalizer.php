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

use App\Model\Translation\Intl\IntlMessageVariation;
use App\Model\Translation\Intl\IntlVariableMessage;

/**
 * Class IntlMessageVariationNormalizer
 */
class IntlVariableMessageNormalizer
{
    /**
     * @param IntlVariableMessage $variable
     *
     * @return IntlVariableMessage
     */
    public function normalize(IntlVariableMessage $variable): object
    {
        $variations = [];
        foreach ($variable as $variation) {
            $selector = $variation->getSelector();
            if ($variable->isPlural()) {
                $selector = ltrim($variation->getSelector(), '=');
            }

            if (is_numeric($selector)) {
                if (is_float($selector)) {
                    $selector = (float) $selector;
                } elseif (is_int($selector)) {
                    $selector = (int) $selector;
                }
            }

            $variations[] = new IntlMessageVariation($selector, $variation->getContent());
        }

        return new IntlVariableMessage($variable->getParam(), $variable->getType(), $variations);
    }
}
