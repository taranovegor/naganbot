<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

use App\Exception\CatchableExceptionInterface;

/**
 * Class GameException
 */
class GameException extends \Exception implements CatchableExceptionInterface
{
    protected $message = 'game.general';

    /**
     * @inheritDoc
     */
    public function isProcessed(): bool
    {
        return true;
    }
}
