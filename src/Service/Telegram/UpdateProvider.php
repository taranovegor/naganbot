<?php
/**
 * Copyright (C) 14.08.20 Egor Taranov
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

namespace App\Service\Telegram;

use App\Exception\Telegram\UpdateCannotBeProvidedException;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpFoundation\RequestStack;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\InvalidJsonException;
use TelegramBot\Api\Types\Update;

/**
 * Class UpdateProvider
 */
class UpdateProvider
{
    private RequestStack $requestStack;

    /**
     * UpdateProvider constructor.
     *
     * @param RequestStack $requestStack
     */
    public function __construct(RequestStack $requestStack)
    {
        $this->requestStack = $requestStack;
    }

    /**
     * @return Update
     *
     * @throws UpdateCannotBeProvidedException
     */
    public function provide(): Update
    {
        $request = $this->requestStack->getMasterRequest();
        if (!$request instanceof Request) {
            throw new UpdateCannotBeProvidedException();
        }

        $content = $this->requestStack->getMasterRequest()->getContent();
        if (!is_string($content)) {
            throw new UpdateCannotBeProvidedException();
        }

        try {
            return Update::fromResponse(BotApi::jsonValidate($content, true));
        } catch (InvalidJsonException $e) {
            throw new UpdateCannotBeProvidedException();
        }
    }
}
