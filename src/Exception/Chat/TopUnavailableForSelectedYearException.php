<?php
/**
 * This file is part of the Nagan Bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Exception\Chat;

use App\Exception\ExceptionMessageInterface;

/**
 * Class TopUnavailableForSelectedYearException
 */
class TopUnavailableForSelectedYearException extends \Exception implements ExceptionMessageInterface
{
    protected $message = 'chat.top_unavailable_for_selected_year';
}
