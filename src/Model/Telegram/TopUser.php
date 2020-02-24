<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Model\Telegram;

use App\Entity\Telegram\User;

/**
 * Class TopUser
 */
class TopUser
{
    /**
     * @var User
     */
    private User $user;

    /**
     * @var int
     */
    private int $numberOfWins;

    /**
     * TopUser constructor.
     *
     * @param User $user
     * @param int  $numberOfWins
     */
    public function __construct(User $user, int $numberOfWins)
    {
        $this->user = $user;
        $this->numberOfWins = $numberOfWins;
    }

    /**
     * @return User
     */
    public function getUser(): User
    {
        return $this->user;
    }

    /**
     * @return int
     */
    public function getNumberOfWins(): int
    {
        return $this->numberOfWins;
    }
}
