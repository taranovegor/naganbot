<?php
/**
 * Copyright (C) 20.09.20 Egor Taranov
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

use App\Exception\Translation\MessagePatternIsInvalidException;
use App\Model\Translation\Intl\IntlMessageVariation;
use App\Model\Translation\Intl\IntlVariableMessage;
use App\Service\Translation\Intl\IntlMessageVariationsProvider;
use App\Service\Translation\Intl\IntlVariableMessageNormalizer;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Cache\Adapter\FilesystemAdapter;
use Symfony\Component\Translation\MessageCatalogueInterface;
use Symfony\Component\Translation\Translator;
use Symfony\Contracts\Cache\ItemInterface;

/**
 * Class IntlMessageVariationsProviderTest
 */
class IntlMessageVariationsProviderTest extends TestCase
{
    public function testProvideWithInvalidId()
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
                    ->willReturn(false)
                ;

                return $catalogue;
            })
        ;

        /** @var MockObject|FilesystemAdapter $cache */
        $cache = $this->createMock(FilesystemAdapter::class);
        /** @var MockObject|IntlVariableMessageNormalizer $normalizer */
        $normalizer = $this->createMock(IntlVariableMessageNormalizer::class);

        $provider = new IntlMessageVariationsProvider($translator, $cache, $normalizer);

        $this->assertEquals(
            null,
            $provider->provide('invalid_message_id')
        );
    }

    public function testProvideWithValidIdAndWithoutArguments()
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
                            other {
                                {step, plural,
                                    =0      {}
                                    other   {}
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

        $provider = new IntlMessageVariationsProvider($translator, $cache, new IntlVariableMessageNormalizer());

        $this->assertEquals(
            new IntlVariableMessage('variation', 'select', [
                new IntlMessageVariation('other', new IntlVariableMessage('step', 'plural', [
                    new IntlMessageVariation(0, ''),
                    new IntlMessageVariation('other', ''),
                ])),
            ]),
            $provider->provide('valid_message_id')
        );
    }

    public function testProvideWithValidIdAndArgument()
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
                            1 {
                                {step, plural,
                                    =0      {0}
                                    =1      {1}
                                    other   {other}
                                }
                            }
                            other {
                                {step, plural,
                                    other   {}
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

        $provider = new IntlMessageVariationsProvider($translator, $cache, new IntlVariableMessageNormalizer());

        $this->assertEquals(
            new IntlVariableMessage('variation', 'select', [
                new IntlMessageVariation(1, new IntlVariableMessage('step', 'plural', [
                    new IntlMessageVariation(0, '0'),
                    new IntlMessageVariation(1, '1'),
                    new IntlMessageVariation('other', 'other'),
                ])),
                new IntlMessageVariation('other', new IntlVariableMessage('step', 'plural', [
                    new IntlMessageVariation('other', ''),
                ])),
            ]),
            $provider->provide('valid_message_id')
        );
    }

    public function testProvideWithInvalidMessagePattern()
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
                        {variation, plural,
                            0 {
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

        $provider = new IntlMessageVariationsProvider($translator, $cache, new IntlVariableMessageNormalizer());

        $this->expectException(MessagePatternIsInvalidException::class);
        $provider->provide('valid_message_id');
    }

    public function testProvideWithInvalidMessageFormat()
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
                        {variation, none,
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

        $provider = new IntlMessageVariationsProvider($translator, $cache, new IntlVariableMessageNormalizer());

        $this->assertEquals(null, $provider->provide('valid_message_id'));
    }
}
