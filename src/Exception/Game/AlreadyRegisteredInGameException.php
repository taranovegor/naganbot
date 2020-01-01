<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

/**
 * Class AlreadyAtGamingTableException
 */
class AlreadyRegisteredInGameException extends GameException
{
    protected $message = 'game.already_at_gaming_table';
}
