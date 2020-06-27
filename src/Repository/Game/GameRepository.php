<?php
/**
 * Copyright (C) 14.08.20 Egor Taranov
 * This file is part of Nagan bot <https://github.com/taranovegor/nagan-bot>.
 *
 * Nagan bot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Nagan bot is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Nagan bot.  If not, see <http://www.gnu.org/licenses/>.
 */

namespace App\Repository\Game;

use App\Entity\Game\Game;
use App\Entity\Telegram\Chat;
use App\Exception\Common\EntityNotFoundException;
use DateTime;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\ORMException;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class GameRepository
 *
 * @method Game|null find($id, $lockMode = null, $lockVersion = null)
 * @method Game|null findOneBy(array $criteria, array $orderBy = null)
 * @method Game[]    findAll()
 * @method Game[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
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
     * @throws ORMException
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
    public function getActiveOrCreatedTodayByChat(Chat $chat): Game
    {
        $from = (new DateTime())->setTime(0, 0, 0);
        $to = (new DateTime())->setTime(23, 59, 59);

        ($qb = $this->createQueryBuilder('gt'))
            ->where($qb->expr()->eq('gt.chat', ':chat'))
            ->andWhere($qb->expr()->orX(
                $qb->expr()->between('gt.createdAt', ':time_from', ':time_to'),
                $qb->expr()->isNull('gt.playedAt')
            ))
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

    /**
     * @param Chat $chat
     *
     * @return Game|null
     */
    public function findLatestByChat(Chat $chat): ?Game
    {
        return $this->findOneBy([
            'chat' => $chat,
        ], ['createdAt' => 'DESC']);
    }

    /**
     * @param Chat $chat
     *
     * @return Game|null
     */
    public function findLatestActiveByChat(Chat $chat): ?Game
    {
        ($qb = $this->createQueryBuilder('g'))
            ->where($qb->expr()->eq('g.chat', ':chat'))
            ->andWhere($qb->expr()->isNull('g.playedAt'))
            ->orderBy('g.createdAt', 'DESC')
            ->setMaxResults(1)
            ->setParameters([
                'chat' => $chat,
            ])
        ;

        return $qb->getQuery()->getOneOrNullResult();
    }
}
