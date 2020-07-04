<?php
/**
 * This file is part of the nagan-bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Exception;

/**
 * Interface ExceptionPayloadInterface
 */
interface TranslatableExceptionInterface
{
    /**
     * @return array
     */
    public function getTranslatableParameters(): array;
}
