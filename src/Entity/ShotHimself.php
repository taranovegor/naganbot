<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Entity;

use Doctrine\ORM\Mapping as ORM;

/**
 * Class ShotHimself
 *
 * @ORM\Table(name="shot_himself")
 * @ORM\Entity(repositoryClass="App\Repository\ShotHimselfRepository")
 */
class ShotHimself
{
    /**
     * @var Gunslinger
     *
     * @ORM\Id()
     * @ORM\ManyToOne(targetEntity="App\Entity\Gunslinger", cascade={"persist"})
     * @ORM\JoinColumn(name="gunslinger_id", referencedColumnName="id", nullable=false)
     */
    private Gunslinger $gunslinger;

    /**
     * ShotHimself constructor.
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
