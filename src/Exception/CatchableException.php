<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception;

use Throwable;

/**
 * Class CatchableException
 */
class CatchableException extends \Exception implements CatchableExceptionInterface
{
    /**
     * @var bool
     */
    private bool $processed;

    /**
     * CatchableException constructor.
     *
     * @param string         $message
     * @param bool           $processed
     * @param int            $code
     * @param Throwable|null $previous
     */
    public function __construct(string $message = "", bool $processed = false, int $code = 0, Throwable $previous = null)
    {
        parent::__construct($message, $code, $previous);
        $this->processed = $processed;
    }


    /**
     * @inheritDoc
     */
    public function isProcessed(): bool
    {
        return $this->processed;
    }
}
