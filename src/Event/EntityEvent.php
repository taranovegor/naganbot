<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;

use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class EntityEvent
 */
abstract class EntityEvent extends Event
{
    /**
     * @var object
     */
    private object $entity;

    /**
     * EntityEvent constructor.
     *
     * @param object $entity
     */
    public function __construct(object $entity)
    {
        $this->entity = $entity;
    }

    /**
     * @return object
     */
    public function getEntity(): object
    {
        return $this->entity;
    }
}
