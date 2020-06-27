<?php
/**
 * Copyright (C) 14.08.20 Egor Taranov
 * This file is part of Nagan bot <https://github.com/taranovegor/nagan-bot>.
 *
 * Nagan bot is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Nagan bot is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Nagan bot.  If not, see <http://www.gnu.org/licenses/>.
 */

namespace App\Service\Common;

use Psr\Container\ContainerInterface;
use Psr\Log\LoggerInterface;

/**
 * Class LoggerFactory
 */
class LoggerFactory
{
    private LoggerInterface $logger;

    private ContainerInterface $container;

    /**
     * LoggerFactory constructor.
     *
     * @param LoggerInterface    $logger
     * @param ContainerInterface $container
     */
    public function __construct(LoggerInterface $logger, ContainerInterface $container)
    {
        $this->logger = $logger;
        $this->container = $container;
    }

    /**
     * @param string $channel
     *
     * @return LoggerInterface
     */
    public function create(string $channel): Loggerinterface
    {
        $loggerId = sprintf('monolog.logger.%s', $channel);
        if (!$this->container->has($loggerId)) {
            return $this->logger;
        }

        return $this->container->get($loggerId);
    }
}
