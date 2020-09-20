<?php
/**
 * Copyright (C) 26.09.20 Egor Taranov
 * This file is part of Nagan bot <https://github.com/taranovegor/nagan-bot>.
 *
 * Nagan bot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Nagan bot is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Nagan bot.  If not, see <http://www.gnu.org/licenses/>.
 */

namespace Service\Translation;

use App\Service\Translation\Intl\IntlMessageVariationsProvider;
use App\Service\Translation\Intl\IntlMessageVariationsRandomizer;
use App\Service\Translation\Intl\IntlVariableMessageNormalizer;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\FilesystemAdapter;
use Symfony\Component\Translation\MessageCatalogueInterface;
use Symfony\Component\Translation\Translator;
use Symfony\Contracts\Cache\ItemInterface;

/**
 * Class IntlMessageVariationsRandomizerTest
 */
class IntlMessageVariationsRandomizerTest extends TestCase
{
    public function testRandWithValidParams()
    {
        /** @var MockObject|Translator $translator */
        $translator = $this->createMock(Translator::class);
        $translator
            ->method('getCatalogue')
            ->willReturnCallback(function () {
                /** @var MockObject|MessageCatalogueInterface $catalogue */
                $catalogue = $this->createMock(MessageCatalogueInterface::class);
                $catalogue
                    ->method('has')
                    ->willReturn(true)
                ;
                $catalogue
                    ->method('get')
                    ->willReturn(<<< TRANSLATION
                        {variation, select,
                            0 {
                                {step, plural,
                                    =0      {0}
                                    other   {other}
                                }
                            }
                            other {
                                {step, plural,
                                    =0      {0}
                                    =1      {1}
                                    other   {other}
                                }
                            }
                        }
                    TRANSLATION)
                ;

                return $catalogue;
            })
        ;

        /** @var MockObject|FilesystemAdapter $cache */
        $cache = $this->createMock(FilesystemAdapter::class);
        $cache
            ->method('get')
            ->will($this->returnCallback(function (string $key, callable $callback) {
                /** @var MockObject|ItemInterface $item */
                $item = $this->createMock(ItemInterface::class);

                return $callback($item);
            }))
        ;

        $randomizer = new IntlMessageVariationsRandomizer(
            new IntlMessageVariationsProvider($translator, $cache, new IntlVariableMessageNormalizer())
        );

        $randomized = $randomizer->rand('id', [
            'variation' => null,
        ]);
        $this->assertCount(1, $randomized);
        $this->assertNotNull($randomized['variation'], 'When randomizing only the variation, it is equal to null');
        $this->assertContains($randomized['variation'], ['0', 'other']);

        $randomized = $randomizer->rand('id', [
            'variation' => null,
            'step' => null,
        ]);
        $this->assertcount(2, $randomized);
        $this->assertNotNull($randomized['variation'], 'When randomizing a variation with a step, the variation is equal to null');
        $this->assertNotNull($randomized['step'], 'When randomizing a variation with a step, the step is equal to null');
        $this->assertContains($randomized['variation'], ['0', 'other']);
        if ('0' === $randomized['variation']) {
            $this->assertContains($randomized['step'], ['0', 'other']);
        } elseif ('other' === $randomized['variation']) {
            $this->assertContains($randomized['step'], ['0', '1', 'other']);
        }

        $randomized = $randomizer->rand('id', [
            'variation' => 0,
            'step' => null,
        ]);
        $this->assertEquals(0, $randomized['variation']);
        $this->assertContains($randomized['step'], ['0', 'other']);

        $randomized = $randomizer->rand('id', [
            'variation' => 0,
            'step' => 0,
        ]);
        $this->assertEquals(0, $randomized['variation']);
        $this->assertEquals(0, $randomized['step']);
    }
}
