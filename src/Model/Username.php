<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Model;

use App\Entity\Telegram\User;
use TelegramBot\Api\Types\User as UserType;

/**
 * Class Username
 */
class Username
{
    /**
     * @var int
     */
    private int $userId;

    /**
     * @var string
     */
    private string $firstName;

    /**
     * @var null|string
     */
    private ?string $lastName;

    /**
     * @var null|string
     */
    private ?string $username;

    /**
     * Username constructor.
     *
     * @param int         $userId
     * @param string      $firstName
     * @param null|string $lastName
     * @param null|string $username
     */
    public function __construct(int $userId, string $firstName, ?string $lastName = null, ?string $username = null)
    {
        $this->userId = $userId;
        $this->firstName = $firstName;
        $this->lastName = $lastName;
        $this->username = $username;
    }

    /**
     * @return string
     */
    public function toString(): string
    {
        if (0 < strlen($this->username)) {
            return $this->username;
        }

        return rtrim(sprintf('%s %s', $this->firstName, $this->lastName), ' ');
    }

    /**
     * @return string
     */
    public function __toString(): string
    {
        return $this->toString();
    }

    /**
     * @param User $user
     *
     * @return Username
     */
    public static function fromUser(User $user): Username
    {
        return new self($user->getId(), $user->getFirstName(), $user->getLastName(), $user->getUsername());
    }

    /**
     * @param UserType $user
     *
     * @return Username
     */
    public static function fromUserType(UserType $user): Username
    {
        return new self($user->getId(), $user->getFirstName(), $user->getLastName(), $user->getUsername());
    }
}
