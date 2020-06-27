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
use App\Exception\Game\ActiveNotFoundException;
use App\Exception\Game\AlreadyJoinedToGameException;
use App\Exception\Game\AlreadyPlayedException;
use App\Manager\Game\GameManager;
use App\Manager\Game\GunslingerManager;
use App\Manager\Telegram\ChatManager;
use App\Manager\Telegram\UserManager;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Doctrine\ORM\ORMException;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class JoinCommand
 */
class JoinCommand extends AbstractCommand implements PublicCommandInterface
{
    public const COMMAND = '/rrjoin';

    private TranslatorInterface $translator;

    private UserManager $userManager;

    private ChatManager $chatManager;

    private GameManager $gameManager;

    private GunslingerManager $gunslingerManager;

    /**
     * JoinCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param UserManager         $userManager
     * @param ChatManager         $chatManager
     * @param GameManager         $gameManager
     * @param GunslingerManager   $gunslingerManager
     */
    public function __construct(TranslatorInterface $translator, UserManager $userManager, ChatManager $chatManager, GameManager $gameManager, GunslingerManager $gunslingerManager)
    {
        $this->translator = $translator;
        $this->userManager = $userManager;
        $this->chatManager = $chatManager;
        $this->gameManager = $gameManager;
        $this->gunslingerManager = $gunslingerManager;
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
        return $this->translator->trans('command.join');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @throws ActiveNotFoundException
     * @throws AlreadyPlayedException
     * @throws EntityNotFoundException
     * @throws ORMException
     * @throws AlreadyJoinedToGameException
     */
    public function execute(BotApi $api, Update $update)
    {
        $user = $this->userManager->get($update->getMessage()->getFrom()->getId());
        $chat = $this->chatManager->get($update->getMessage()->getChat()->getId());

        try {
            $game = $this->gameManager->getLatestActiveByChat($chat);
        } catch (EntityNotFoundException $e) {
            throw new ActiveNotFoundException();
        }

        $this->gunslingerManager->create($game, $user);
    }
}
