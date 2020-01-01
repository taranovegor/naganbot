<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

/**
 * Class NotFoundActiveGameException
 */
class NotFoundActiveGameException extends GameException
{
    protected $message = 'game.not_found_active_game';
}
