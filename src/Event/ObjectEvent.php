<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Event;

use Symfony\Contracts\EventDispatcher\Event;

/**
 * Class ObjectEvent
 */
class ObjectEvent extends Event
{
    /**
     * @var object
     */
    private object $object;

    /**
     * ObjectEvent constructor.
     *
     * @param object $object
     */
    public function __construct(object $object)
    {
        $this->object = $object;
    }

    /**
     * @return object
     */
    public function getObject(): object
    {
        return $this->object;
    }
}
