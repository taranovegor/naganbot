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
use App\Exception\Game\AlreadyCreatedException;
use App\Exception\Game\AlreadyJoinedToGameException;
use App\Exception\Game\AlreadyPlayedException;
use App\Manager\Game\GameManager;
use App\Manager\Telegram\ChatManager;
use App\Manager\Telegram\UserManager;
use BoShurik\TelegramBotBundle\Telegram\Command\AbstractCommand;
use BoShurik\TelegramBotBundle\Telegram\Command\PublicCommandInterface;
use Doctrine\ORM\ORMException;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Update;

/**
 * Class StartCommand
 */
class CreateCommand extends AbstractCommand implements PublicCommandInterface
{
    public const NAME = '/ngcreate';

    private TranslatorInterface $translator;

    private ChatManager $chatManager;

    private UserManager $userManager;

    private GameManager $gameManager;

    /**
     * StartCommand constructor.
     *
     * @param TranslatorInterface $translator
     * @param ChatManager         $chatManager
     * @param UserManager         $userManager
     * @param GameManager         $gameManager
     */
    public function __construct(TranslatorInterface $translator, ChatManager $chatManager, UserManager $userManager, GameManager $gameManager)
    {
        $this->translator = $translator;
        $this->chatManager = $chatManager;
        $this->userManager = $userManager;
        $this->gameManager = $gameManager;
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
        return $this->translator->trans('command.create');
    }

    /**
     * @param BotApi $api
     * @param Update $update
     *
     * @throws AlreadyCreatedException
     * @throws AlreadyJoinedToGameException
     * @throws AlreadyPlayedException
     * @throws EntityNotFoundException
     * @throws ORMException
     */
    public function execute(BotApi $api, Update $update)
    {
        $chat = $this->chatManager->get($update->getMessage()->getChat()->getId());
        $user = $this->userManager->get($update->getMessage()->getFrom()->getId());

        $this->gameManager->create($chat, $user);
    }
}
