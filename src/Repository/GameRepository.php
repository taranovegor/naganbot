<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository;

use App\Entity\Game;
use App\Entity\Telegram\Chat;
use App\Exception\EntityNotFoundException;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;
use Ramsey\Uuid\Doctrine\UuidBinaryType;

/**
 * Class GameRepository
 */
class GameRepository extends ServiceEntityRepository
{
    /**
     * GameRepository constructor.
     *
     * @param ManagerRegistry $registry
     */
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, Game::class);
    }

    /**
     * @param Game $game
     *
     * @throws \Doctrine\ORM\ORMException
     */
    public function add(Game $game)
    {
        $this->getEntityManager()->persist($game);
    }

    /**
     * @param Chat $chat
     *
     * @return Game
     *
     * @throws EntityNotFoundException
     */
    public function getByChat(Chat $chat): Game
    {
        $object = $this->findOneBy([
            'chat' => $chat,
        ], ['createdAt' => 'DESC']);

        if (!$object instanceof Game) {
            throw new EntityNotFoundException();
        }

        return $object;
    }

    /**
     * @param Chat $chat
     *
     * @return Game
     *
     * @throws EntityNotFoundException
     */
    public function getCreatedTodayByChat(Chat $chat): Game
    {
        $from = (new \DateTime())->setTime(0, 0, 0);
        $to = (new \DateTime())->setTime(23, 59, 59);

        ($qb = $this->createQueryBuilder('gt'))
            ->where($qb->expr()->eq('gt.chat', ':chat'))
            ->andWhere($qb->expr()->between('gt.createdAt', ':time_from', ':time_to'))
            ->setParameters([
                'chat' => $chat->getId(),
                'time_from' => $from,
                'time_to' => $to,
            ])
            ->setFirstResult(0)
            ->setMaxResults(1)
        ;

        $object = $qb->getQuery()->getResult()[0] ?? null;
        if (!$object instanceof Game) {
            throw new EntityNotFoundException();
        }

        return $object;
    }
}
