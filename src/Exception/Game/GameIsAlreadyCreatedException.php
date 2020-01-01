<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

/**
 * Class GameTableAlreadyCreatedException
 */
class GameIsAlreadyCreatedException extends GameException
{
    protected $message = 'game.game_table_already_created';
}
