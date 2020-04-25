<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Entity\Telegram;

use App\Repository\Telegram\UserRepository;
use Doctrine\ORM\Mapping as ORM;
use TelegramBot\Api\Types\User as UserType;

/**
 * Class User
 *
 * @ORM\Table(name="telegram_user")
 * @ORM\Entity(repositoryClass=UserRepository::class)
 */
class User
{
    /**
     * @var int
     *
     * @ORM\Id()
     * @ORM\Column(name="id", type="integer", unique=true, nullable=false)
     */
    private int $id;

    /**
     * @var bool
     *
     * @ORM\Column(name="is_bot", type="boolean", nullable=false)
     */
    private bool $isBot;

    /**
     * @var string
     *
     * @ORM\Column(name="first_name", type="string", nullable=false)
     */
    private string $firstName;

    /**
     * @var null|string
     *
     * @ORM\Column(name="last_name", type="string", nullable=true)
     */
    private ?string $lastName;

    /**
     * @var null|string
     *
     * @ORM\Column(name="username", type="string", nullable=true)
     */
    private ?string $username;

    /**
     * @var null|string
     *
     * @ORM\Column(name="language_code", type="string", nullable=true)
     */
    private ?string $languageCode;

    /**
     * User constructor.
     *
     * @param int    $id
     * @param string $firstName
     */
    public function __construct(int $id, string $firstName)
    {
        $this->id = $id;
        $this->firstName = $firstName;
        $this->lastName = null;
        $this->username = null;
        $this->languageCode = null;
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }

    /**
     * @return bool
     */
    public function isBot(): bool
    {
        return $this->isBot;
    }

    /**
     * @param bool $isBot
     *
     * @return User
     */
    public function setIsBot(bool $isBot): User
    {
        $this->isBot = $isBot;

        return $this;
    }

    /**
     * @return string
     */
    public function getFirstName(): string
    {
        return $this->firstName;
    }

    /**
     * @param string $firstName
     *
     * @return User
     */
    public function setFirstName(string $firstName): User
    {
        $this->firstName = $firstName;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getLastName(): ?string
    {
        return $this->lastName;
    }

    /**
     * @param null|string $lastName
     *
     * @return User
     */
    public function setLastName(?string $lastName): User
    {
        $this->lastName = $lastName;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getUsername(): ?string
    {
        return $this->username;
    }

    /**
     * @param null|string $username
     *
     * @return User
     */
    public function setUsername(?string $username): User
    {
        $this->username = $username;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getLanguageCode(): ?string
    {
        return $this->languageCode;
    }

    /**
     * @param null|string $languageCode
     *
     * @return User
     */
    public function setLanguageCode(?string $languageCode): User
    {
        $this->languageCode = $languageCode;

        return $this;
    }

    /**
     * @param User $user
     *
     * @return bool
     */
    public function isSame(User $user): bool
    {
        return $this->getId() === $user->getId();
    }

    /**
     * @param UserType $user
     *
     * @return User
     */
    public function updateFromUserType(UserType $user): User
    {
        return $this
            ->setIsBot($user->isBot())
            ->setFirstName($user->getFirstName())
            ->setLastName($user->getLastName())
            ->setUsername($user->getUsername())
            ->setLanguageCode($user->getLanguageCode())
        ;
    }

    /**
     * @param UserType $user
     *
     * @return User
     */
    public static function createFromUserType(UserType $user): User
    {
        return (new self($user->getId(), $user->getFirstName()))
            ->setIsBot($user->isBot())
            ->setFirstName($user->getFirstName())
            ->setLastName($user->getLastName())
            ->setUsername($user->getUsername())
            ->setLanguageCode($user->getLanguageCode())
        ;
    }
}
