<?php
/**
 * This file is part of the nagan-bot application.
 *
 * For the full copyright and license information, please view the LICENSE file that was distributed with this source code.
 */

namespace App\Manager;

use App\Entity\Telegram\Chat;
use App\Exception\EntityNotFoundException;
use App\Repository\Telegram\ChatRepository;

/**
 * Class ChatManager
 */
class ChatManager
{
    private ChatRepository $repository;

    /**
     * ChatManager constructor.
     *
     * @param ChatRepository $repository
     */
    public function __construct(ChatRepository $repository)
    {
        $this->repository = $repository;
    }

    /**
     * @param int $id
     *
     * @return Chat
     *
     * @throws EntityNotFoundException
     */
    public function get(int $id): Chat
    {
        $object = $this->repository->find($id);
        if (!$object instanceof Chat) {
            throw new EntityNotFoundException();
        }

        return $object;
    }
}
