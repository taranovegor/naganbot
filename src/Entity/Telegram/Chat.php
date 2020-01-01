<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Entity\Telegram;

use Doctrine\ORM\Mapping as ORM;
use TelegramBot\Api\Types\Chat as ChatType;

/**
 * Class Chat
 *
 * @ORM\Table(name="telegram_chat")
 * @ORM\Entity(repositoryClass="App\Repository\Telegram\ChatRepository")
 */
class Chat
{
    /**
     * @var int
     *
     * @ORM\Id()
     * @ORM\Column(name="id", type="bigint", unique=true, nullable=false)
     */
    private int $id;

    /**
     * @var string
     *
     * @ORM\Column(name="type", type="string", nullable=false)
     */
    private string $type;

    /**
     * @var null|string
     *
     * @ORM\Column(name="title", type="string", nullable=true)
     */
    private ?string $title;

    /**
     * @var null|string
     *
     * @ORM\Column(name="first_name", type="string", nullable=true)
     */
    private ?string $firstName;

    /**
     * @var null|string
     *
     * @ORM\Column(name="last_name", type="string", nullable=true)
     */
    private ?string $lastName;

    /**
     * @var null|string
     *
     * @ORM\Column(name="invite_link", type="string", nullable=true)
     */
    private ?string $inviteLink;

    /**
     * Chat constructor.
     *
     * @param int    $id
     * @param string $type
     */
    public function __construct(int $id, string $type)
    {
        $this->id = $id;
        $this->type = $type;
        $this->title = null;
        $this->firstName = null;
        $this->lastName = null;
        $this->inviteLink = null;
    }

    /**
     * @return int
     */
    public function getId(): int
    {
        return $this->id;
    }


    /**
     * @return string
     */
    public function getType(): string
    {
        return $this->type;
    }

    /**
     * @param string $type
     *
     * @return Chat
     */
    public function setType(string $type): Chat
    {
        $this->type = $type;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getTitle(): ?string
    {
        return $this->title;
    }

    /**
     * @param null|string $title
     *
     * @return Chat
     */
    public function setTitle(?string $title): Chat
    {
        $this->title = $title;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getFirstName(): ?string
    {
        return $this->firstName;
    }

    /**
     * @param null|string $firstName
     *
     * @return Chat
     */
    public function setFirstName(?string $firstName): Chat
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
     * @return Chat
     */
    public function setLastName(?string $lastName): Chat
    {
        $this->lastName = $lastName;

        return $this;
    }

    /**
     * @return null|string
     */
    public function getInviteLink(): ?string
    {
        return $this->inviteLink;
    }

    /**
     * @param null|string $inviteLink
     *
     * @return Chat
     */
    public function setInviteLink(?string $inviteLink): Chat
    {
        $this->inviteLink = $inviteLink;

        return $this;
    }

    /**
     * @param ChatType $chat
     *
     * @return Chat
     */
    public function updateFromChatType(ChatType $chat): Chat
    {
        return $this
            ->setTitle($chat->getTitle())
            ->setFirstName($chat->getFirstName())
            ->setLastName($chat->getLastName())
            ->setInviteLink($chat->getInviteLink())
        ;
    }

    /**
     * @param Chat $chat
     *
     * @return bool
     */
    public function isSame(Chat $chat): bool
    {
        return $chat->getId() === $this->getId();
    }

    /**
     * @param ChatType $chat
     *
     * @return Chat
     */
    public static function createFromChatType(ChatType $chat): Chat
    {
        return (new self($chat->getId(), $chat->getType()))
            ->setTitle($chat->getTitle())
            ->setFirstName($chat->getFirstName())
            ->setLastName($chat->getLastName())
            ->setInviteLink($chat->getInviteLink())
        ;
    }
}
