<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Exception\Game;

use App\Exception\TranslatableExceptionInterface;

/**
 * Class NeedToLieDownException
 */
class GameIsAlreadyPlayedException extends GameException implements TranslatableExceptionInterface
{
    protected $message = 'game.need_to_lie_down';

    /**
     * @return array
     */
    public function getTranslatableParameters(): array
    {
        return [
            'remaining_time' => (new \DateTime())->diff((new \DateTime('+1 day'))->setTime(0, 0, 0)),
        ];
    }
}
