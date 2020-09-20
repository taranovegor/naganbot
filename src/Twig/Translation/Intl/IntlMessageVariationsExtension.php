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

namespace App\Twig\Translation\Intl;

use App\Constant\Common\Translation\Domain;
use App\Exception\Translation\MessagePatternIsInvalidException;
use App\Service\Translation\Intl\IntlMessageVariationsRandomizer;
use Psr\Cache\InvalidArgumentException;
use Twig\Extension\AbstractExtension;
use Twig\TwigFunction;

/**
 * Class IntlMessageVariationsExtension
 */
class IntlMessageVariationsExtension extends AbstractExtension
{
    private IntlMessageVariationsRandomizer $variationsRandomizer;

    /**
     * IntlMessageVariationsExtension constructor.
     *
     * @param IntlMessageVariationsRandomizer $variationsRandomizer
     */
    public function __construct(IntlMessageVariationsRandomizer $variationsRandomizer)
    {
        $this->variationsRandomizer = $variationsRandomizer;
    }

    /**
     * @inheritDoc
     */
    public function getFunctions()
    {
        return [
            new TwigFunction('intl_rand_variation_params', [$this, 'intlRandVariationParams']),
        ];
    }

    /**
     * @param string      $id
     * @param array       $parameters
     * @param string|null $domain
     * @param string|null $locale
     *
     * @return array
     *
     * @throws InvalidArgumentException
     * @throws MessagePatternIsInvalidException
     */
    public function intlRandVariationParams(string $id, array $parameters = [], string $domain = Domain::DEFAULT, string $locale = null): array
    {
        return $this->variationsRandomizer->rand($id, $parameters, $domain, $locale);
    }
}
