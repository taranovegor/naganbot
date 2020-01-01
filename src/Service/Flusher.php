<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Service;

use Doctrine\ORM\EntityManagerInterface;

/**
 * Class Flusher
 */
class Flusher
{
    /**
     * @var EntityManagerInterface
     */
    private EntityManagerInterface $em;

    /**
     * Flusher constructor.
     *
     * @param EntityManagerInterface $em
     */
    public function __construct(EntityManagerInterface $em)
    {
        $this->em = $em;
    }

    /**
     * @return void
     */
    public function flush(): void
    {
        $this->em->flush();
    }
}
