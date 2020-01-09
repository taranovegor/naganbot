<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception;

/**
 * Interface ExceptionInterface
 */
interface CatchableExceptionInterface
{
    /**
     * @return bool
     */
    public function isProcessed(): bool;

    /**
     * @return mixed
     */
    public function getMessage();
}
