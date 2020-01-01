<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\MessageBuilder;

use App\Model\Message;

/**
 * Class AbstractMessageBuilder
 */
abstract class AbstractMessageBuilder
{
    /**
     * @var Message
     */
    protected Message $message;

    /**
     * AbstractMessageBuilder constructor.
     */
    public function __construct()
    {
        $this->message = new Message();
    }
}
