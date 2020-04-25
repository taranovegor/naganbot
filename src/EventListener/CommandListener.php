<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventListener;

use App\Exception\CatchableException;
use App\Exception\ChatOnlyException;
use App\Service\TelegramUpdateHandler;
use BoShurik\TelegramBotBundle\Event\UpdateEvent;

/**
 * Class CommandListener
 */
class CommandListener
{
    /**
     * @var TelegramUpdateHandler
     */
    private TelegramUpdateHandler $dataSaver;

    /**
     * CommandSubscriber constructor.
     *
     * @param TelegramUpdateHandler $dataSaver
     */
    public function __construct(TelegramUpdateHandler $dataSaver)
    {
        $this->dataSaver = $dataSaver;
    }

    /**
     * @param UpdateEvent $event
     *
     * @throws CatchableException
     * @throws ChatOnlyException
     */
    public function onUpdate(UpdateEvent $event)
    {
        $this->dataSaver->handle($event->getUpdate());
    }
}
