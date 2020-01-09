<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Repository\Telegram;

use App\Entity\Telegram\Chat;
use App\Exception\EntityNotFoundException;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class ChatRepository
 */
class ChatRepository extends ServiceEntityRepository
{
    /**
     * ChatRepository constructor.
     *
     * @param ManagerRegistry $registry
     */
    public function __construct(ManagerRegistry $registry)
    {
        parent::__construct($registry, Chat::class);
    }

    /**
     * @param Chat $chat
     *
     * @throws \Doctrine\ORM\ORMException
     */
    public function add(Chat $chat)
    {
        $this->getEntityManager()->persist($chat);
    }

    /**
     * @param int $id
     *
     * @return Chat
     *
     * @throws EntityNotFoundException
     */
    public function get(int $id): Chat
    {
        $object = $this->find($id);

        if (!$object instanceof Chat) {
            throw new EntityNotFoundException();
        }

        return $object;
    }
}
