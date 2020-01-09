<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository\Telegram;

use App\Entity\Telegram\User;
use App\Exception\EntityNotFoundException;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class UserRepository
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
}
