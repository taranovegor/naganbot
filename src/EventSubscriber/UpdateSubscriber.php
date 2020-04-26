<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventSubscriber;

use App\Exception\CatchableException;
use App\Exception\ChatOnlyException;
use App\Service\TelegramUpdateHandler;
use BoShurik\TelegramBotBundle\Event\UpdateEvent;
use Psr\Log\LoggerInterface;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;

/**
 * Class UpdateSubscriber
 */
class UpdateSubscriber implements EventSubscriberInterface
{
    /**
     * @var TelegramUpdateHandler
     */
    private TelegramUpdateHandler $dataSaver;

    /**
     * @var LoggerInterface
     */
    private LoggerInterface $logger;

    /**
     * CommandSubscriber constructor.
     *
     * @param TelegramUpdateHandler $dataSaver
     * @param LoggerInterface       $logger
     */
    public function __construct(TelegramUpdateHandler $dataSaver, LoggerInterface $logger)
    {
        $this->dataSaver = $dataSaver;
        $this->logger = $logger;
    }

    /**
     * @inheritDoc
     */
    public static function getSubscribedEvents()
    {
        return [
            UpdateEvent::class => [
                ['onUpdate', 255],
            ],
        ];
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
