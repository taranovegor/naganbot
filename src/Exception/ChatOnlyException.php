<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception;

/**
 * Class ChatOnlyException
 */
class ChatOnlyException extends \Exception implements CatchableExceptionInterface
{
    protected $message = 'general.chat_only';

    /**
     * @inheritDoc
     */
    public function isProcessed(): bool
    {
        return true;
    }
}