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

namespace App\Repository\Telegram;

use App\Entity\Game\Game;
use App\Entity\Game\Gunslinger;
use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\NonUniqueResultException;
use Doctrine\ORM\NoResultException;
use Doctrine\ORM\ORMException;
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
     * @throws ORMException
     */
    public function add(User $user)
    {
        $this->getEntityManager()->persist($user);
    }

    /**
     * @param User $user
     * @param Chat $chat
     *
     * @return int
     *
     * @throws NoResultException
     * @throws NonUniqueResultException
     */
    public function numberOfWinsForUserInChat(User $user, Chat $chat): int
    {
        ($qb = $this->createQueryBuilder('_user'))
            ->select($qb->expr()->count('_user.id'))
            ->innerJoin(
                Gunslinger::class,
                '_gunslinger',
                Join::WITH,
                $qb->expr()->eq('_user.id', '_gunslinger.user')
            )
            ->innerJoin(
                Game::class,
                '_game',
                Join::WITH,
                $qb->expr()->eq('_gunslinger.game', '_game.id')
            )
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

        return (int) $qb->getQuery()->getSingleScalarResult();
    }
}
