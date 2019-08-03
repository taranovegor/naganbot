<?php

namespace App\Commands;

abstract class Command
{
    /**
     * @var \Formapro\TelegramBot\Update
     */
    private $update;
    
    /**
     * Command constructor.
     * @param \Formapro\TelegramBot\Update $update
     */
    public function __construct(\Formapro\TelegramBot\Update $update)
    {
        $this->update = $update;
    }
    
    /**
     * @return object|string
     */
    public abstract function handle();
}