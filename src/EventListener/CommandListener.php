<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventListener;

use App\Exception\CatchableException;
use App\Service\TelegramDataSaver;
use BoShurik\TelegramBotBundle\Event\Telegram\UpdateEvent;

/**
 * Class CommandSubscriber
 */
class CommandListener
{
    /**
     * @var TelegramDataSaver
     */
    private TelegramDataSaver $dataSaver;

    /**
     * CommandSubscriber constructor.
     *
     * @param TelegramDataSaver $dataSaver
     */
    public function __construct(TelegramDataSaver $dataSaver)
    {
        $this->dataSaver = $dataSaver;
    }

    /**
     * @param UpdateEvent $event
     *
     * @throws CatchableException
     * @throws \App\Exception\ChatOnlyException
     */
    public function onUpdate(UpdateEvent $event)
    {
        $this->dataSaver->handle($event->getUpdate());
    }
}
