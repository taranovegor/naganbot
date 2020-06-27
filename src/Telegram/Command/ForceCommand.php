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

use App\Exception\Common\EntityNotFoundException;
use App\Exception\Game\AlreadyPlayedException;
use App\Exception\Game\FailedToShuffleArrayException;
use App\Exception\Game\GunslingerNotFoundException;
use App\Exception\Game\NotEnoughGunslingersException;
use App\Manager\Game\GameManager;
use App\Manager\Telegram\ChatManager;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class ForceCommand
 */
class ForceCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrforce';

    private TranslatorInterface $translator;

    private GameManager $gameManager;

    private ChatManager $chatManager;

    /**
     * ForceCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param GameManager         $gameManager
     * @param ChatManager         $chatManager
     */
    public function __construct(TranslatorInterface $translator, GameManager $gameManager, ChatManager $chatManager)
    {
        $this->translator = $translator;
        $this->gameManager = $gameManager;
        $this->chatManager = $chatManager;
    }

    /**
     * @inheritDoc
     */
    public function getName()
    {
        return self::COMMAND;
    }

    /**
     * @inheritDoc
     */
    public function getDescription()
    {
        return $this->translator->trans('command.force');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @throws GunslingerNotFoundException
     * @throws EntityNotFoundException
     * @throws FailedToShuffleArrayException
     */
    public function execute(BotApi $api, Update $update)
    {
        try {
            $chat = $this->chatManager->get($update->getMessage()->getChat()->getId());
            $game = $this->gameManager->getLatestByChat($chat);
            $this->gameManager->play($game);
        } catch (NotEnoughGunslingersException | AlreadyPlayedException $e) {
        }
    }
}
