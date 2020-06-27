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

namespace App\EventSubscriber\Kernel;

use App\Constant\Telegram\ParseMode;
use App\Exception\ExceptionMessageInterface;
use App\Exception\ParameterizedExceptionMessageInterface;
use App\Exception\Telegram\UpdateCannotBeProvidedException;
use App\Service\DateTime\DateTimeDiffFormatter;
use App\Service\Telegram\UpdateProvider;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\ExceptionEvent;
use Symfony\Component\HttpKernel\Exception\HttpException;
use Symfony\Component\HttpKernel\KernelEvents;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Exception;
use TelegramBot\Api\InvalidArgumentException;

/**
 * Class ExceptionSubscriber
 */
class ExceptionSubscriber implements EventSubscriberInterface
{
    private UpdateProvider $updateProvider;

    private TranslatorInterface $translator;

    private DateTimeDiffFormatter $dateTimeDiffFormatter;

    private BotApi $api;

    /**
     * ExceptionSubscriber constructor.
     *
     * @param UpdateProvider        $updateProvider
     * @param TranslatorInterface   $translator
     * @param DateTimeDiffFormatter $dateTimeDiffFormatter
     * @param BotApi                $api
     */
    public function __construct(UpdateProvider $updateProvider, TranslatorInterface $translator, DateTimeDiffFormatter $dateTimeDiffFormatter, BotApi $api)
    {
        $this->updateProvider = $updateProvider;
        $this->translator = $translator;
        $this->dateTimeDiffFormatter = $dateTimeDiffFormatter;
        $this->api = $api;
    }

    /**
     * {@inheritDoc}
     */
    public static function getSubscribedEvents()
    {
        return [
            KernelEvents::EXCEPTION => [
                ['sendMessageToChat', 0],
            ],
        ];
    }

    /**
     * @param ExceptionEvent $event
     *
     * @throws UpdateCannotBeProvidedException
     * @throws Exception
     * @throws InvalidArgumentException
     */
    public function sendMessageToChat(ExceptionEvent $event): void
    {
        $throwable = $event->getThrowable();
        if (!$throwable instanceof ExceptionMessageInterface) {
            return;
        }

        if ($throwable instanceof ParameterizedExceptionMessageInterface) {
            $parameters = [];
            foreach ($throwable->getMessageParameters() as $key => $value) {
                if ($value instanceof \DateInterval) {
                    $value = $this->dateTimeDiffFormatter->format($value);
                }

                $parameters[sprintf('%%%s%%', $key)] = $value;
            }
        } else {
            $parameters = [];
        }

        $this->api->sendMessage(
            $this->updateProvider->provide()->getMessage()->getChat()->getId(),
            $this->translator->trans($throwable->getMessage(), $parameters, 'errors'),
            ParseMode::MARKDOWN
        );

        $event->setThrowable(new HttpException(Response::HTTP_NO_CONTENT));
    }
}
