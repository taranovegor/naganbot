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
use App\Telegram\Command\JoinCommand;
use Iterator;
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

    /**
     * GameMessageBuilder constructor.
     *
     * @param Environment $twig
     */
    public function __construct(Environment $twig)
    {
        $this->twig = $twig;
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
        return $this->twig->render('Game/create.md.twig', [
            'command' => JoinCommand::COMMAND,
        ]);
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
        return $this->twig->render('Game/joined.md.twig', [
            'game' => $game,
        ]);
    }

    /**
     * @return string[]|Iterator
     *
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function buildPlay(): Iterator
    {
        $variation = rand(0, 2);

        for ($step = 0; $step <= 1; $step++) {
            yield $this->twig->render('Game/play.md.twig', [
                'variation' => $variation,
                'step' => $step,
            ]);
        }
    }
}
