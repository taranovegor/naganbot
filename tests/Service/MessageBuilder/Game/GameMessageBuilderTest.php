<?php
/**
 * This file is part of the Nagan bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace Service\MessageBuilder\Game;

use App\Service\MessageBuilder\Game\GameMessageBuilder;
use App\Service\Translation\Intl\IntlMessageVariationsProvider;
use App\Service\Translation\Intl\IntlMessageVariationsRandomizer;
use PHPUnit\Framework\MockObject\MockObject;
use Symfony\Bundle\FrameworkBundle\Test\WebTestCase;

/**
 * Class GameMessageBuilderTest
 */
class GameMessageBuilderTest extends WebTestCase
{
    public function testBuildReadyToPlay()
    {
        self::bootKernel();

        /** @var IntlMessageVariationsProvider $provider */
        $provider = self::$container->get(IntlMessageVariationsProvider::class);
        /** @var MockObject|IntlMessageVariationsRandomizer $randomizer */
        $randomizer = $this->createMock(IntlMessageVariationsRandomizer::class);

        $builder = new GameMessageBuilder(
            self::$container->get('twig'),
            $provider,
            $randomizer
        );

        foreach ($provider->provide('game.ready_to_play') as $variation) {
            $randomizer
                ->method('rand')
                ->willReturn(['variation' => $variation->getSelector()])
            ;

            $built = $builder->buildReadyToPlay();
            $this->assertIsIterable($built, 'Built message not iterable');
            $steps = 0;
            $prevStep = null;
            foreach ($built as $step) {
                $this->assertNotEmpty($step, 'Step message is empty');
                $this->assertIsString($step, 'Step message not string');
                $this->assertNotTrue(
                    null !== $prevStep && $step !== $prevStep,
                    'Current step message equals to previous'
                );
                ++$steps;
            }
            $this->assertGreaterThan(0, $steps, 'invalid number of steps');
        }
    }
}
