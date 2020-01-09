<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Model\Telegram\Type;

/**
 * Class Chat
 */
class Chat
{
    public const TYPE_PRIVATE = 'private';
    public const TYPE_GROUP = 'group';
    public const TYPE_SUPERGROUP = 'supergroup';
    public const TYPE_CHANNEL = 'channel';
}
