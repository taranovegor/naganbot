<?php
/**
 * This file is part of the Nagan bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Twig\Strategy\Escaper;

use Twig\Environment;

/**
 * Interface EscaperStrategyInterface
 */
interface EscaperStrategyInterface
{
    /**
     * @param Environment $twig
     * @param string      $string
     * @param             $charset
     *
     * @return string
     */
    public function __invoke(Environment $twig, string $string, $charset): string;

    /**
     * @return string
     */
    public static function getName(): string;
}
