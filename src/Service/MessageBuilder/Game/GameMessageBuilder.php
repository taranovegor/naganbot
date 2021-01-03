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

namespace App\Service\MessageBuilder\Game;

use App\Entity\Game\Game;
use App\Exception\Translation\MessagePatternIsInvalidException;
use App\Service\Translation\Intl\IntlMessageVariationsProvider;
use App\Service\Translation\Intl\IntlMessageVariationsRandomizer;
use Iterator;
use Psr\Cache\InvalidArgumentException;
use Twig\Environment;
use Twig\Error\LoaderError;
use Twig\Error\RuntimeError;
use Twig\Error\SyntaxError;

/**
 * Class GameMessageBuilder
 */
class GameMessageBuilder
{
    private Environment $twig;

    private IntlMessageVariationsProvider $variationsProvider;

    private IntlMessageVariationsRandomizer $variationsRandomizer;

    /**
     * GameMessageBuilder constructor.
     *
     * @param Environment                     $twig
     * @param IntlMessageVariationsProvider   $variationsProvider
     * @param IntlMessageVariationsRandomizer $variationsRandomizer
     */
    public function __construct(Environment $twig, IntlMessageVariationsProvider $variationsProvider, IntlMessageVariationsRandomizer $variationsRandomizer)
    {
        $this->twig = $twig;
        $this->variationsProvider = $variationsProvider;
        $this->variationsRandomizer = $variationsRandomizer;
    }

    /**
     * @return string
     *
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function buildCreate(): string
    {
        return $this->twig->render('Game/Game/create.md.twig');
    }

    /**
     * @param Game $game
     *
     * @return string
     *
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function buildJoined(Game $game): string
    {
        return $this->twig->render('Game/Game/joined.md.twig', [
            'game' => $game,
        ]);
    }

    /**
     * @return string[]|Iterator
     *
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws MessagePatternIsInvalidException
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function buildReadyToPlay(): Iterator
    {
        $variation = $this->variationsProvider->provide('game.ready_to_play')->getVariation(
            $this->variationsRandomizer->rand(
                'game.play',
                ['variation' => null]
            )['variation']
        );

        foreach ($variation->getContent() as $step) {
            yield $this->twig->render('Game/Game/ready_to_play.md.twig', [
                'variation' => $variation->getSelector(),
                'step' => $step->getSelector(),
            ]);
        }
    }

    /**
     * @return string
     *
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function buildPlayedWithNuclearBullet(): string
    {
        return $this->twig->render('Game/Game/played_with_nuclear_bullet.md.twig');
    }
}
