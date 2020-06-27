<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Entity\Game;

use App\Entity\Game\Gunslinger;
use App\Entity\Telegram\Chat;
use App\Entity\Telegram\User;
use DateTime;
use Doctrine\Common\Collections\ArrayCollection;
use Doctrine\Common\Collections\Collection;
use Doctrine\Common\Collections\Criteria;
use Doctrine\ORM\Mapping as ORM;
use Ramsey\Uuid\Uuid;
use Ramsey\Uuid\UuidInterface;
use App\Repository\Game\GameRepository;

/**
 * Class GameTable
 *
 * @ORM\Table(name="game_game")
 * @ORM\Entity(repositoryClass=GameRepository::class)
 */
class Game
{
    /**
     * @var UuidInterface
     *
     * @ORM\Id()
     * @ORM\Column(name="id", type="uuid_binary", unique=true, nullable=false)
     */
    private UuidInterface $id;

    /**
     * @var Chat
     *
     * @ORM\ManyToOne(targetEntity=Chat::class, cascade={"persist"})
     * @ORM\JoinColumn(name="chat_id", referencedColumnName="id", nullable=false)
     */
    private Chat $chat;

    /**
     * @var User
     *
     * @ORM\ManyToOne(targetEntity=User::class, cascade={"persist"})
     * @ORM\JoinColumn(name="owner_user_id", referencedColumnName="id", nullable=false)
     */
    private User $owner;

    /**
     * @var Gunslinger[]|Collection
     *
     * @ORM\OneToMany(targetEntity=Gunslinger::class, cascade={"persist"}, mappedBy="game")
     * @ORM\OrderBy({"joinedAt" = "ASC"})
     */
    private Collection $gunslingers;

    /**
     * @var DateTime
     *
     * @ORM\Column(name="created_at", type="datetime", nullable=false)
     */
    private DateTime $createdAt;

    /**
     * @var null|DateTime
     *
     * @ORM\Column(name="played_out_at", type="datetime", nullable=true)
     */
    private ?DateTime $playedAt;

    /**
     * GameTable constructor.
     *
     * @param Chat $chat
     * @param User $owner
     */
    public function __construct(Chat $chat, User $owner)
    {
        $this->id = Uuid::uuid4();
        $this->chat = $chat;
        $this->owner = $owner;
        $this->gunslingers = new ArrayCollection();
        $this->createdAt = new DateTime();
        $this->playedAt = null;
    }

    /**
     * @return UuidInterface
     */
    public function getId(): UuidInterface
    {
        return $this->id;
    }

    /**
     * @return Chat
     */
    public function getChat(): Chat
    {
        return $this->chat;
    }

    /**
     * @return User
     */
    public function getOwner(): User
    {
        return $this->owner;
    }

    /**
     * @return Gunslinger[]|Collection
     */
    public function getGunslingers(): Collection
    {
        return $this->gunslingers;
    }

    /**
     * @param Gunslinger $gunslinger
     *
     * @return Game
     */
    public function addGunslinger(Gunslinger $gunslinger): Game
    {
        if (!$this->gunslingers->contains($gunslinger)) {
            $this->gunslingers->add($gunslinger);
        }

        return $this;
    }

    /**
     * @param User $user
     *
     * @return null|Gunslinger
     */
    public function getGunslingerByUser(User $user): ?Gunslinger
    {
        /** @var Gunslinger $gunslinger */
        foreach ($this->gunslingers->toArray() as $gunslinger) {
            if ($user->isSame($gunslinger->getUser())) {
                return $gunslinger;
            }
        }

        return null;
    }

    /**
     * @return DateTime
     */
    public function getCreatedAt(): DateTime
    {
        return $this->createdAt;
    }

    /**
     * @return DateTime|null
     */
    public function getPlayedAt(): ?DateTime
    {
        return $this->playedAt;
    }

    /**
     * @return bool
     */
    public function isPlayed(): bool
    {
        return $this->playedAt instanceof DateTime;
    }

    /**
     * @return Game
     */
    public function setAsPlayed(): Game
    {
        if (null === $this->playedAt) {
            $this->playedAt = new DateTime();
        }

        return $this;
    }

    /**
     * @return bool
     */
    public function isCreatedToday(): bool
    {
        return 0 === (int) (clone $this
                ->createdAt
                ->setTime(0, 0, 0)
            )
            ->diff(new DateTime())
            ->format('%R%a')
        ;
    }
}
