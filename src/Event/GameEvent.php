<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;

use App\Entity\Game\Game;
use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class GameEvent
 */
class GameEvent extends Event
{
    public const CREATED = 'event.game.created';
    public const PLAYED = 'event.game.played';

    /**
     * @var Game
     */
    private Game $game;

    /**
     * GameEvent constructor.
     *
     * @param Game $game
     */
    public function __construct(Game $game)
    {
        $this->game = $game;
    }

    /**
     * @return Game
     */
    public function getGame(): Game
    {
        return $this->game;
    }
}
