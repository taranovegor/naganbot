<?php
/**
 * This file is part of the Nagan bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Twig\Strategy\Escaper;

use Twig\Environment;

/**
 * Class MarkdownEscaperStrategy
 */
class MarkdownEscaperStrategy implements EscaperStrategyInterface
{
    /**
     * @inheritDoc
     */
    public function __invoke(Environment $twig, string $string, $charset): string
    {
        return addcslashes($string, '\`*{}[]()#+-_.!');
    }

    /**
     * @inheritDoc
     */
    public static function getName(): string
    {
        return 'markdown';
    }
}
