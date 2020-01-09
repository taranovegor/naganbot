<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;


use App\Entity\Game;
use App\Entity\Gunslinger;

/**
 * Class GunslingerJoinedEvent
 */
class GunslingerJoinedEvent extends GameEvent
{
    public const EVENT = self::GUNSLINGER_JOINED;

    /**
     * @var Gunslinger
     */
    private Gunslinger $gunslinger;

    /**
     * GunslingerJoinedEvent constructor.
     *
     * @param Game       $gameTable
     * @param Gunslinger $gunslinger
     */
    public function __construct(Game $gameTable, Gunslinger $gunslinger)
    {
        parent::__construct($gameTable);
        $this->gunslinger = $gunslinger;
    }

    /**
     * @return Gunslinger
     */
    public function getGunslinger(): Gunslinger
    {
        return $this->gunslinger;
    }
}
