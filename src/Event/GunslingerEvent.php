<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;

use App\Entity\Game\Gunslinger;
use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class GunslingerEvent
 */
class GunslingerEvent extends Event
{
    public const JOINED = 'event.gunslinger.joined';
    public const SHOT_HIMSELF = 'event.gunslinger.shot_himself';

    /**
     * @var Gunslinger
     */
    private Gunslinger $gunslinger;

    /**
     * GunslingerEvent constructor.
     *
     * @param Gunslinger $gunslinger
     */
    public function __construct(Gunslinger $gunslinger)
    {
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
