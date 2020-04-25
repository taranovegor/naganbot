<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Model\Telegram\Type;

/**
 * Class KeyboardButton
 *
 * @see https://core.telegram.org/bots/api#keyboardbutton
 */
class KeyboardButton
{
    /**
     * @var string
     */
    private string $text;

    /**ffsd
     * @var bool
     */
    private bool $requestContact;

    /**
     * @var bool
     */
    private bool $requestLocation;

    /**
     * KeyboardButton constructor.
     *
     * @param string $text
     * @param bool   $requestContact
     * @param bool   $requestLocation
     */
    public function __construct(string $text, bool $requestContact, bool $requestLocation)
    {
        $this->text = $text;
        $this->requestContact = $requestContact;
        $this->requestLocation = $requestLocation;
    }

    /**
     * @param string $text
     * @param bool   $requestContact
     * @param bool   $requestLocation
     *
     * @return self
     */
    public static function create(string $text, bool $requestContact = false, bool $requestLocation = false): self
    {
        return new self($text, $requestContact, $requestLocation);
    }

    /**
     * @param string $text
     * @param bool   $requestContact
     * @param bool   $requestLocation
     *
     * @return array
     */
    public static function createMarkdown(string $text, bool $requestContact = false, bool $requestLocation = false): array
    {
        return (new self($text, $requestContact, $requestLocation))->toArray();
    }

    /**
     * @return array
     */
    public function toArray()
    {
        return [
            'text' => $this->text,
            'request_contact' => $this->requestContact,
            'request_location' => $this->requestContact,
        ];
    }
}
