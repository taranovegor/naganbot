<?php
/**
 * This file is part of the Nagan bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\DependencyInjection\Configurator\Twig;

use App\Twig\Strategy\Escaper\EscaperStrategyInterface;
use Symfony\Bundle\TwigBundle\DependencyInjection\Configurator\EnvironmentConfigurator as BaseEnvironmentConfigurator;
use Twig\Environment;
use Twig\Extension\EscaperExtension;

/**
 * Class EnvironmentConfigurator
 */
class EnvironmentConfigurator
{
    private BaseEnvironmentConfigurator $configurator;

    private iterable $strategiesIterators;

    /**
     * EnvironmentConfigurator constructor.
     *
     * @param BaseEnvironmentConfigurator $configurator
     * @param iterable                    $strategiesGenerators
     */
    public function __construct(BaseEnvironmentConfigurator $configurator, iterable $strategiesGenerators)
    {
        $this->configurator = $configurator;
        $this->strategiesIterators = $strategiesGenerators;
    }

    /**
     * @param Environment $twig
     *
     * @throws \Exception
     */
    public function configure(Environment $twig)
    {
        $this->configurator->configure($twig);

        /** @var EscaperExtension $escaperExtension */
        $escaperExtension = $twig->getExtension(EscaperExtension::class);

        foreach ($this->strategiesIterators as $strategiesIterator) {
            /** @var EscaperStrategyInterface $strategy */
            foreach ($strategiesIterator as $strategy) {
                $escaperExtension->setEscaper($strategy::getName(), $strategy);
            }
        }
    }
}
