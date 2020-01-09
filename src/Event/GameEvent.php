<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;

use App\Entity\Game;
use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class GameEvent
 */
class GameEvent extends Event
{
    public const CREATED = 'game.created';
    public const GUNSLINGER_JOINED = 'game.gunslinger_joined';
    public const REQUIRED_NUMBER_OF_GUNSLINGER_REACHED = 'game.required_number_of_gunslinger_reached';
    public const GUNSLINGER_SHOT_HIMSELF = 'game.gunslinger_shot_himself';

    /**
     * @var Game
     */
    private Game $gameTable;

    /**
     * GameEvent constructor.
     *
     * @param Game $gameTable
     */
    public function __construct(Game $gameTable)
    {
        $this->gameTable = $gameTable;
    }

    /**
     * @return Game
     */
    public function getGameTable(): Game
    {
        return $this->gameTable;
    }
}
