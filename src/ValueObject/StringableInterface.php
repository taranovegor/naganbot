<?php
/**
 * This file is part of the WOG application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\ValueObject;

/**
 * Interface StringableInterface
 */
interface StringableInterface
{
    /**
     * @return string
     */
    public function __toString(): string;
}
