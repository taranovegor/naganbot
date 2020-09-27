<?php
/**
 * Copyright (C) 27.09.20 Egor Taranov
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

namespace Telegram\Command;

use App\Entity\Game\Game;
use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use App\Exception\Common\EntityNotFoundException;
use App\Exception\Game\ActiveNotFoundException;
use App\Manager\Game\GameManager;
use App\Manager\Game\GunslingerManager;
use App\Manager\Telegram\ChatManager;
use App\Manager\Telegram\UserManager;
use App\Telegram\Command\JoinCommand;
use PHPUnit\Framework\MockObject\MockObject;
use PHPUnit\Framework\TestCase;
use Symfony\Component\Translation\Translator;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\Types\Message;
use TelegramBot\Api\Types\Update;

/**
 * Class JoinCommandTest
 */
class JoinCommandTest extends TestCase
{
    public function testExecuteWithoutActiveGame()
    {
        /** @var MockObject|UserManager $userManager */
        $userManager = $this->createMock(UserManager::class);
        $userManager
            ->method('get')
            ->willReturn($this->createMock(User::class))
        ;

        /** @var MockObject|ChatManager $chatManager */
        $chatManager = $this->createMock(ChatManager::class);
        $chatManager
            ->method('get')
            ->willReturn($this->createMock(Chat::class))
        ;

        /** @var MockObject|GameManager $gameManager */
        $gameManager = $this->createMock(GameManager::class);
        $gameManager
            ->method('getActiveOrCreatedTodayByChat')
            ->willThrowException(new EntityNotFoundException())
        ;

        $command = new JoinCommand(
            $this->createMock(Translator::class),
            $userManager,
            $chatManager,
            $gameManager,
            $this->createMock(GunslingerManager::class)
        );

        $this->expectException(ActiveNotFoundException::class);
        $command->execute(
            $this->createMock(BotApi::class),
            $this->createUpdateMock()
        );
    }

    public function testExecuteWithActiveGame()
    {
        /** @var MockObject|UserManager $userManager */
        $userManager = $this->createMock(UserManager::class);
        $userManager
            ->method('get')
            ->willReturn($this->createMock(User::class))
        ;

        /** @var MockObject|ChatManager $chatManager */
        $chatManager = $this->createMock(ChatManager::class);
        $chatManager
            ->method('get')
            ->willReturn($this->createMock(Chat::class))
        ;

        /** @var MockObject|GameManager $gameManager */
        $gameManager = $this->createMock(GameManager::class);
        $gameManager
            ->method('getActiveOrCreatedTodayByChat')
            ->willReturn($this->createMock(Game::class))
        ;

        $command = new JoinCommand(
            $this->createMock(Translator::class),
            $userManager,
            $chatManager,
            $gameManager,
            $this->createMock(GunslingerManager::class)
        );

        $command->execute(
            $this->createMock(BotApi::class),
            $this->createUpdateMock()
        );
        $this->assertTrue(true);
    }

    /**
     * @return MockObject|Update
     */
    protected function createUpdateMock()
    {
        /** @var MockObject|User $user */
        $user = $this->createMock(User::class);
        $user
            ->method('getId')
            ->willReturn(0)
        ;

        /** @var MockObject|Chat $chat */
        $chat = $this->createMock(Chat::class);
        $chat
            ->method('getId')
            ->willReturn(0)
        ;

        /** @var MockObject|Message $message */
        $message = $this->createMock(Message::class);
        $message
            ->method('getFrom')
            ->willReturn($user)
        ;
        $message
            ->method('getChat')
            ->willReturn($chat)
        ;

        /** @var MockObject|Update $update */
        $update =  $this->createMock(Update::class);
        $update
            ->method('getMessage')
            ->willReturn($message)
        ;

        return $update;
    }
}
