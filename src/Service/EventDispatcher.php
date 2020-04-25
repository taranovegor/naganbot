<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Service;

use Symfony\Component\EventDispatcher\EventDispatcherInterface;

/**
 * Class EventDispatcher
 */
class EventDispatcher
{
    /**
     * @var EventDispatcherInterface
     */
    private EventDispatcherInterface $eventDispatcher;

    /**
     * EventDispatcher constructor.
     *
     * @param EventDispatcherInterface $eventDispatcher
     */
    public function __construct(EventDispatcherInterface $eventDispatcher)
    {
        $this->eventDispatcher = $eventDispatcher;
    }

    /**
     * @param             $event
     * @param string|null $eventName
     *
     * @return object
     */
    public function dispatch($event, string $eventName = null): object
    {
        return $this->eventDispatcher->dispatch($event, $eventName);
    }
}
