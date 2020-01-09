<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

/**
 * Class ShotDeadNotFoundException
 */
class ShotDeadNotFoundException extends GameException
{
    protected $message = 'game.misfire';
}
