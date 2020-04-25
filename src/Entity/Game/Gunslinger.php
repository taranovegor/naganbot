<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Entity\Game;

use App\Entity\Game\Game;
use App\Entity\Telegram\User;
use DateTime;
use Doctrine\ORM\Mapping as ORM;
use Ramsey\Uuid\Uuid;
use Ramsey\Uuid\UuidInterface;
use App\Repository\Game\GunslingerRepository;

/**
 * Class Gunslinger
 *
 * @ORM\Table(name="game_gunslinger")
 * @ORM\Entity(repositoryClass=GunslingerRepository::class)
 */
class Gunslinger
{
    /**
     * @var UuidInterface
     *
     * @ORM\Id()
     * @ORM\Column(name="id", type="uuid_binary", unique=true, nullable=false)
     */
    private UuidInterface $id;

    /**
     * @var Game
     *
     * @ORM\ManyToOne(targetEntity=Game::class, cascade={"persist"}, inversedBy="gunslingers")
     * @ORM\JoinColumn(name="game_id", referencedColumnName="id", nullable=false)
     */
    private Game $game;

    /**
     * @var User
     *
     * @ORM\ManyToOne(targetEntity=User::class, cascade={"persist"})
     * @ORM\JoinColumn(name="user_id", referencedColumnName="id", nullable=false)
     */
    private User $user;

    /**
     * @var DateTime
     *
     * @ORM\Column(name="joined_at", type="datetime", nullable=false)
     */
    private DateTime $joinedAt;

    /**
     * @var bool
     *
     * @ORM\Column(name="shot_himself", type="boolean", nullable=false)
     */
    private bool $shotHimself = false;

    /**
     * Gunslinger constructor.
     *
     * @param Game $gameTable
     * @param User $user
     */
    public function __construct(Game $gameTable, User $user)
    {
        $this->id = Uuid::uuid4();
        $this->game = $gameTable;
        $this->user = $user;
        $this->joinedAt = new DateTime();
    }

    /**
     * @return UuidInterface
     */
    public function getId(): UuidInterface
    {
        return $this->id;
    }

    /**
     * @return Game
     */
    public function getGame(): Game
    {
        return $this->game;
    }

    /**
     * @return User
     */
    public function getUser(): User
    {
        return $this->user;
    }

    /**
     * @return DateTime
     */
    public function getJoinedAt(): DateTime
    {
        return $this->joinedAt;
    }

    /**
     * @return bool
     */
    public function isShotHimself(): bool
    {
        return $this->shotHimself;
    }

    /**
     * @return Gunslinger
     */
    public function shot(): Gunslinger
    {
        $this->shotHimself = true;

        return $this;
    }

    /**
     * @param Gunslinger $gunslinger
     *
     * @return bool
     */
    public function isSame(Gunslinger $gunslinger): bool
    {
        return $this->getId()->equals($gunslinger->getId());
    }
}
