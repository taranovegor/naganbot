<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Model;

/**
 * Class Parameters
 */
class Parameters
{
    /**
     * @var null|array
     */
    private ?array $matches;

    /**
     * Parameters constructor.
     *
     * @param string $text
     * @param string $parameters must contains regex, for example: (\d)+ (\w)
     * @param string $command
     */
    public function __construct(string $text, string $parameters, string $command = null)
    {
        $pattern = '/';
        if (null !== $command) {
            $pattern .= sprintf('\\%s ', $command);
        }
        $pattern .= sprintf('%s/m', $parameters);
        preg_match($pattern, $text, $matches);
        array_shift($matches);
        $this->matches = $matches;
    }

    /**
     * @return mixed
     */
    public function getFirst()
    {
        return $this->matches[0] ?? null;
    }

    /**
     * @param \TelegramBot\Api\Types\Message $message
     * @param string                         $parameters must contains regex, for example: (\d)+ (\w)
     * @param string|null                    $command
     *
     * @return Parameters
     */
    public static function fromMessageType(\TelegramBot\Api\Types\Message $message, string $parameters, string $command = null): Parameters
    {
        return new self($message->getText(), $parameters, $command);
    }
}
