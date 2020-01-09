<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

/**
 * Class NeedToLieDownException
 */
class GameIsAlreadyPlayedException extends GameException
{
    protected $message = 'game.need_to_lie_down';
}
