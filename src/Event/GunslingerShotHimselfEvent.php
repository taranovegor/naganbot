<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;

use App\Entity\Game;
use App\Entity\ShotHimself;

/**
 * Class GunslingerShotHimselfEvent
 */
class GunslingerShotHimselfEvent extends GameEvent
{
    public const EVENT = self::GUNSLINGER_SHOT_HIMSELF;

    /**
     * @var ShotHimself
     */
    private ShotHimself $shotHimself;

    /**
     * GunslingerShotHimselfEvent constructor.
     *
     * @param Game        $gameTable
     * @param ShotHimself $shotHimself
     */
    public function __construct(Game $gameTable, ShotHimself $shotHimself)
    {
        parent::__construct($gameTable);
        $this->shotHimself = $shotHimself;
    }

    /**
     * @return ShotHimself
     */
    public function getShotHimself(): ShotHimself
    {
        return $this->shotHimself;
    }
}
