<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository\Telegram;

use App\Entity\Game\Game;
use App\Entity\Game\Gunslinger;
use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Exception\EntityNotFoundException;
use App\Model\Telegram\TopUser;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\Query\Expr\Join;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class UserRepository
 *
 * @method User|null find($id, $lockMode = null, $lockVersion = null)
 * @method User|null findOneBy(array $criteria, array $orderBy = null)
 * @method User[]    findAll()
 * @method User[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
 */
class UserRepository extends ServiceEntityRepository
{
    /**
     * UserRepository constructor.
     *
     * @param ManagerRegistry $registry
     */
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, User::class);
    }

    /**
     * @param User $user
     *
     * @throws \Doctrine\ORM\ORMException
     */
    public function add(User $user)
    {
        $this->getEntityManager()->persist($user);
    }

    /**
     * @param int $id
     *
     * @return User
     *
     * @throws EntityNotFoundException
     */
    public function get(int $id): User
    {
        $object = $this->find($id);

        if (!$object instanceof User) {
            throw new EntityNotFoundException();
        }

        return $object;
    }

    /**
     * @param Chat $chat
     * @param int  $limit
     *
     * @return array|TopUser[]
     */
    public function topByChat(Chat $chat, int $limit = 5): array
    {
        ($qb = $this->createQueryBuilder('_user'))
            ->select('_user AS user, COUNT(_user.id) AS number_of_wins')
            ->innerJoin(Gunslinger::class, '_gunslinger', Join::WITH, $qb->expr()->eq('_user.id', '_gunslinger.user'))
            ->innerJoin(Game::class, '_game', Join::WITH, '_gunslinger.game = _game.id')
            ->where($qb->expr()->eq('_game.chat', ':chat'))
            ->andWhere($qb->expr()->eq('_gunslinger.shotHimself', ':shot_himself'))
            ->andWhere($qb->expr()->isNotNull('_game.playedAt'))
            ->groupBy('_user.id')
            ->orderBy('number_of_wins', 'DESC')
            ->setMaxResults($limit)
            ->setParameters([
                'chat' => $chat,
                'shot_himself' => true,
            ])
        ;

        $results = [];
        foreach ($qb->getQuery()->getResult() as $result) {
            $results[] = new TopUser($result['user'], $result['number_of_wins']);
        }

        return $results;
    }

    /**
     * @param User $user
     * @param Chat $chat
     *
     * @return TopUser
     */
    public function getStatisticByUserInChat(User $user, Chat $chat): TopUser
    {
        ($qb = $this->createQueryBuilder('_user'))
            ->select($qb->expr()->count('_user.id'))
            ->innerJoin(Gunslinger::class, '_gunslinger', Join::WITH, $qb->expr()->eq('_user.id', '_gunslinger.user'))
            ->innerJoin(Game::class, '_game', Join::WITH, '_gunslinger.game = _game.id')
            ->where($qb->expr()->eq('_user', ':user'))
            ->andWhere($qb->expr()->eq('_game.chat', ':chat'))
            ->andWhere($qb->expr()->eq('_gunslinger.shotHimself', ':shot_himself'))
            ->andWhere($qb->expr()->isNotNull('_game.playedAt'))
            ->setParameters([
                'user' => $user,
                'chat' => $chat,
                'shot_himself' => true,
            ])
        ;

        return new TopUser($user, $qb->getQuery()->getSingleScalarResult());
    }
}
