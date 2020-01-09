<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception;

use Exception;

/**
 * Class EntityNotFoundException
 */
class EntityNotFoundException extends Exception implements CatchableExceptionInterface
{
    protected $message = 'general.entity_not_found';

    /**
     * @inheritDoc
     */
    public function isProcessed(): bool
    {
        return true;
    }
}
