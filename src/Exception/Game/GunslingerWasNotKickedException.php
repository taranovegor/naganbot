<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

/**
 * Class WoundWasNotFatalException
 */
class GunslingerWasNotKickedException extends GameException
{
    protected $message = 'game.wound_was_not_fatal';
}
