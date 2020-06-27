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
use App\Exception\Common\EntityNotFoundException;
use Doctrine\Bundle\DoctrineBundle\Repository\ServiceEntityRepository;
use Doctrine\ORM\ORMException;
use Doctrine\ORM\Query\Expr\Join;
use Doctrine\Persistence\ManagerRegistry;

/**
 * Class ChatRepository
 *
 * @method Chat|null find($id, $lockMode = null, $lockVersion = null)
 * @method Chat|null findOneBy(array $criteria, array $orderBy = null)
 * @method Chat[]    findAll()
 * @method Chat[]    findBy(array $criteria, array $orderBy = null, $limit = null, $offset = null)
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
     * @throws ORMException
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

    /**
     * @param Chat $chat
     * @param int  $limit
     *
     * @return array
     */
    public function numberOfWinsForEachChatMemberInChat(Chat $chat, int $limit = 10): array
    {
        ($qb = $this->createQueryBuilder('_chat'))
            ->select(sprintf('_user AS user, %s AS number_of_wins', $qb->expr()->count('_user.id')))
            ->innerJoin(
                Game::class,
                '_game',
                Join::WITH,
                $qb->expr()->eq('_game.chat', '_chat')
            )
            ->innerJoin(
                Gunslinger::class,
                '_gunslinger',
                Join::WITH,
                $qb->expr()->eq('_gunslinger.game', '_game')
            )
            ->innerJoin(
                User::class,
                '_user',
                Join::WITH,
                $qb->expr()->eq('_user', '_gunslinger.user')
            )
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

        return $qb->getQuery()->getResult();
    }
}
