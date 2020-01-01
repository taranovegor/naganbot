<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Model;

/**
 * Class Message
 */
class Message
{
    /**
     * @var string
     */
    private string $message;

    /**
     * Message constructor.
     */
    public function __construct()
    {
        $this->clear();
    }

    /**
     * @return Message
     */
    public function clear(): Message
    {
        $this->message = '';

        return $this;
    }

    /**
     * @param string $str
     *
     * @return Message
     */
    public function add(string $str): Message
    {
        $this->message .= $str;

        return $this;
    }

    /**
     * @return Message
     */
    public function nextLine(): Message
    {
        $this->message .= PHP_EOL;

        return $this;
    }

    /**
     * @param string $str
     *
     * @return Message
     */
    public function addLine(string $str): Message
    {
        return $this->add($str)->nextLine();
    }

    /**
     * @return Message
     */
    public function addSpace(): Message
    {
        return $this->add(' ');
    }

    /**
     * @return string
     */
    public function toString(): string
    {
        return rtrim($this->message, PHP_EOL);
    }

    /**
     * @return string
     */
    public function __toString(): string
    {
        return $this->toString();
    }
}
