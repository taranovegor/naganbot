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

namespace App\Telegram\Command;

use App\Constant\Telegram\ParseMode;
use App\Exception\Common\EntityNotFoundException;
use App\Manager\Game\GameManager;
use App\Manager\Telegram\ChatManager;
use App\Service\MessageBuilder\Game\GameMessageBuilder;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Exception;
use TelegramBot\Api\InvalidArgumentException;
use TelegramBot\Api\Types\Update;
use Twig\Error\LoaderError;
use Twig\Error\RuntimeError;
use Twig\Error\SyntaxError;

/**
 * Class JoinedCommand
 */
class JoinedCommand extends AbstractCommand implements PublicCommandInterface
{
    public const NAME = '/ngjoined';

    private TranslatorInterface $translator;

    private GameManager $gameManager;

    private ChatManager $chatManager;

    private GameMessageBuilder $messageBuilder;

    /**
     * JoinedCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param GameManager         $gameManager
     * @param ChatManager         $chatManager
     * @param GameMessageBuilder  $messageBuilder
     */
    public function __construct(TranslatorInterface $translator, GameManager $gameManager, ChatManager $chatManager, GameMessageBuilder $messageBuilder)
    {
        $this->translator = $translator;
        $this->gameManager = $gameManager;
        $this->chatManager = $chatManager;
        $this->messageBuilder = $messageBuilder;
    }

    /**
     * @inheritDoc
     */
    public function getName()
    {
        return self::NAME;
    }

    /**
     * @inheritDoc
     */
    public function getDescription()
    {
        return $this->translator->trans('command.joined');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @throws Exception
     * @throws InvalidArgumentException
     * @throws LoaderError
     * @throws RuntimeError
     * @throws SyntaxError
     */
    public function execute(BotApi $api, Update $update)
    {
        try {
            $chat = $this->chatManager->get($update->getMessage()->getChat()->getId());
            $game = $this->gameManager->getLatestByChat($chat);
        } catch (EntityNotFoundException $e) {
            return;
        }

        $api->sendMessage(
            $game->getChat()->getId(),
            $this->messageBuilder->buildJoined($game),
            ParseMode::MARKDOWN
        );
    }
}
