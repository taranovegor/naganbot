<?php
/**
 * This file is part of the Nagan bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Service\Game;

use App\Entity\Game\Game;

/**
 * Class NuclearBulletChecker
 */
class NuclearBulletChecker
{
    private int $min;

    private int $max;

    /**
     * NuclearBulletChecker constructor.
     *
     * @param int $min
     * @param int $max
     */
    public function __construct(int $min, int $max)
    {
        $this->min = $min;
        $this->max = $max;
    }

    /**
     * @param Game $game
     *
     * @return bool
     */
    public function isNuclearBullet(Game $game): bool
    {
        return false !== strpos($game->getId()->toString(), pack('H*', rand($this->min, $this->max)));
    }
}
