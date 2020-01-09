<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository;

use App\Entity\Gunslinger;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class GunslingerRepository
 */
class GunslingerRepository extends ServiceEntityRepository
{
    /**
     * GunslingerRepository constructor.
     *
     * @param ManagerRegistry $registry
     */
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, Gunslinger::class);
    }

    /**
     * @param Gunslinger $gunslinger
     *
     * @throws \Doctrine\ORM\ORMException
     */
    public function add(Gunslinger $gunslinger)
    {
        $this->getEntityManager()->persist($gunslinger);
    }
}
