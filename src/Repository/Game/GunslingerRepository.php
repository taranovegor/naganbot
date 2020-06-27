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
use App\Entity\Game\Gunslinger;
use App\Entity\Telegram\User;
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

    /**
     * @param Gunslinger $gunslinger
     *
     * @throws ORMException
     */
    public function add(Gunslinger $gunslinger): void
    {
        $this->getEntityManager()->persist($gunslinger);
    }

    /**
     * @param User $user
     * @param Game $game
     *
     * @return Gunslinger|null
     */
    public function findByUserInGame(User $user, Game $game): ?Gunslinger
    {
        return $this->findOneBy([
            'game' => $game,
            'user' => $user,
        ]);
    }

    /**
     * @param Game $game
     *
     * @return Gunslinger[]
     */
    public function findByGame(Game $game): array
    {
        return $this->findBy([
            'game' => $game,
        ], ['joinedAt' => 'DESC']);
    }
}
