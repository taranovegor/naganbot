<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository;

use App\Entity\Game;
use App\Entity\Gunslinger;
use App\Entity\ShotHimself;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;
use Ramsey\Uuid\Doctrine\UuidBinaryType;

/**
 * Class ShotHimselfRepository
 */
class ShotHimselfRepository extends ServiceEntityRepository
{
    /**
     * ShotHimselfRepository constructor.
     *
     * @param ManagerRegistry $registry
     */
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, ShotHimself::class);
    }

    /**
     * @param ShotHimself $shotHimself
     *
     * @throws \Doctrine\ORM\ORMException
     */
    public function add(ShotHimself $shotHimself)
    {
        $this->getEntityManager()->persist($shotHimself);
    }

    /**
     * @param Game $game
     *
     * @return null|ShotHimself
     */
    public function findByGameTable(Game $game): ?ShotHimself
    {
        ($qb = $this->createQueryBuilder('sh'))
            ->innerJoin(Gunslinger::class, 'gl')
            ->where($qb->expr()->eq('gl.game', ':game'))
            ->andWhere($qb->expr()->eq('sh.gunslinger', 'gl.id'))
            ->setParameter('game', $game->getId(), UuidBinaryType::NAME)
        ;

        return $qb->getQuery()->getResult()[0] ?? null;
    }
}
