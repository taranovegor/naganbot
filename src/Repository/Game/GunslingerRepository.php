<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository\Game;

use App\Entity\Game\Gunslinger;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\ORMException;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class GunslingerRepository
 *
 * @method Gunslinger|null find($id, $lockMode = null, $lockVersion = null)
 * @method Gunslinger|null findOneBy(array $criteria, array $orderBy = null)
 * @method Gunslinger[]    findAll()
 * @method Gunslinger[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
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
}
