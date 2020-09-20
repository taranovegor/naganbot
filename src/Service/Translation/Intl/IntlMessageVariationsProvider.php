<?php
/*
 * Copyright © 2008 by Yii Software LLC (http://www.yiisoft.com)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 *  * Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in
 *    the documentation and/or other materials provided with the
 *    distribution.
 *  * Neither the name of Yii Software LLC nor the names of its
 *    contributors may be used to endorse or promote products derived
 *    from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
 * BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
 * ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * Originally forked from
 * https://github.com/yiisoft/yii2/blob/2.0.15/framework/i18n/MessageFormatter.php
 */

namespace App\Service\Translation\Intl;

use App\Constant\Common\Translation\Domain;
use App\Exception\Translation\MessageFormatNotSupportedException;
use App\Exception\Translation\MessagePatternIsInvalidException;
use App\Model\Translation\Intl\IntlMessageVariation;
use App\Model\Translation\Intl\IntlVariableMessage;
use Symfony\Component\Cache\Adapter\FilesystemAdapter;
use Symfony\Component\Translation\Translator;
use Symfony\Contracts\Cache\ItemInterface;

/**
 * Class IntlMessageVariationsProvider
 *
 * It only supports the following message formats:
 *  * plural formatting for english ('one' and 'other' selectors)
 *  * select format
 *  * simple parameters
 *
 * It does NOT support the ['apostrophe-friendly' syntax](https://php.net/MessageFormatter.formatMessage).
 * Also messages that are working with the fallback implementation are not necessarily compatible with the
 * PHP intl MessageFormatter so do not rely on the fallback if you are able to install intl extension somehow.
 *
 * @author Alexander Makarov <sam@rmcreative.ru>
 * @author Carsten Brandt <mail@cebe.cc>
 * @author Nicolas Grekas <p@tchwork.com>
 */
class IntlMessageVariationsProvider
{
    private Translator $translator;

    private FilesystemAdapter $cache;

    private IntlVariableMessageNormalizer $normalizer;

    /**
     * IntlMessageVariationsProvider constructor.
     *
     * @param Translator                    $translator
     * @param FilesystemAdapter             $cache
     * @param IntlVariableMessageNormalizer $normalizer
     */
    public function __construct(Translator $translator, FilesystemAdapter $cache, IntlVariableMessageNormalizer $normalizer)
    {
        $this->translator = $translator;
        $this->cache = $cache;
        $this->normalizer = $normalizer;
    }

    /**
     * @param string      $id
     * @param string      $domain
     * @param string|null $locale
     *
     * @return IntlVariableMessage|null
     *
     * @throws MessagePatternIsInvalidException
     * @throws \Psr\Cache\InvalidArgumentException
     */
    public function provide(string $id, string $domain = Domain::DEFAULT, string $locale = null): ?IntlVariableMessage
    {
        $catalogue = $this->translator->getCatalogue($locale);
        if (!$catalogue->has($id, $domain)) {
            return null;
        }
        $pattern = $catalogue->get($id, $domain);

        return $this->cache->get(
            sprintf('translation.intl.intl_message_variations_provider.%s', $id),
            function (ItemInterface $item) use ($pattern) {
                $parsedTokens = $this->parseTokens(
                    self::tokenizePattern($pattern)
                );
                if (!$parsedTokens instanceof IntlVariableMessage) {
                    return null;
                }

                return $parsedTokens;
            }
        );
    }

    /**
     * @param array $tokens
     *
     * @return IntlVariableMessage|array|null
     *
     * @throws MessagePatternIsInvalidException
     */
    private function parseTokens(array $tokens)
    {
        foreach ($tokens as $i => $token) {
            if (\is_array($token)) {
                try {
                    $tokens[$i] = $this->parseToken($token);
                } catch (MessageFormatNotSupportedException $e) {
                    continue;
                }
            }
        }

        return $tokens[1] ?? $tokens[0] ?? null;
    }

    /**
     * Parses pattern based on ICU grammar.
     *
     * @see http://icu-project.org/apiref/icu4c/classMessageFormat.html#details
     *
     * @param array $token
     *
     * @return IntlVariableMessage
     *
     * @throws MessageFormatNotSupportedException
     * @throws MessagePatternIsInvalidException
     */
    private function parseToken(array $token)
    {
        $param = trim($token[0]);
        $type = isset($token[1]) ? trim($token[1]) : 'none';
        switch ($type) {
            case 'date':
            case 'time':
            case 'spellout':
            case 'ordinal':
            case 'duration':
            case 'choice':
            case 'selectordinal':
            case 'number':
            case 'none':
                throw new MessageFormatNotSupportedException(
                    sprintf('"%s" message format not supported', $type)
                );

            case 'plural':
            case 'select':
                if (!isset($token[2])) {
                    throw new MessagePatternIsInvalidException();
                }
                $variation = self::tokenizePattern($token[2]);
                $c = \count($variation);
                $isValid = false;
                $variations = [];
                for ($i = 0; 1 + $i < $c; ++$i) {
                    if (\is_array($variation[$i]) || !\is_array($variation[1 + $i])) {
                        throw new MessagePatternIsInvalidException();
                    }
                    $selector = trim($variation[$i++]);
                    $variations[] = new IntlMessageVariation(
                        $selector,
                        $this->parseTokens(
                            self::tokenizePattern(implode(',', $variation[$i]))
                        )
                    );

                    $isValid = false === $isValid && 'other' === $selector;
                }
                if ($isValid) {
                    return $this->normalizer->normalize(
                        new IntlVariableMessage($param, $type, $variations)
                    );
                }
                break;
        }

        throw new MessagePatternIsInvalidException();
    }

    /**
     * @param string $pattern
     *
     * @return array|string[]
     *
     * @throws MessagePatternIsInvalidException
     */
    private static function tokenizePattern(string $pattern): array
    {
        if (false === $start = $pos = strpos($pattern, '{')) {
            return [$pattern];
        }

        $depth = 1;
        $tokens = [substr($pattern, 0, $pos)];

        while (true) {
            $open = strpos($pattern, '{', 1 + $pos);
            $close = strpos($pattern, '}', 1 + $pos);

            if (false === $open) {
                if (false === $close) {
                    break;
                }
                $open = \strlen($pattern);
            }

            if ($close > $open) {
                ++$depth;
                $pos = $open;
            } else {
                --$depth;
                $pos = $close;
            }

            if (0 === $depth) {
                $tokens[] = explode(
                    ',',
                    substr($pattern, 1 + $start, $pos - $start - 1),
                    3
                );
                $start = 1 + $pos;
                $tokens[] = substr($pattern, $start, $open - $start);
                $start = $open;
            }

            if (0 !== $depth && (false === $open || false === $close)) {
                break;
            }
        }

        if ($depth) {
            throw new MessagePatternIsInvalidException();
        }

        return $tokens;
    }
}
