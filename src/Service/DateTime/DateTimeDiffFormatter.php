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

namespace App\Service\DateTime;

use Symfony\Contracts\Translation\TranslatorInterface;

/**
 * Class DateTimeFormatter
 */
class DateTimeDiffFormatter
{
    private const UNITS = [
        'd' => 'day',
        'h' => 'hour',
        'i' => 'minute',
    ];

    private TranslatorInterface $translator;

    /**
     * DateTimeFormatter constructor.
     *
     * @param TranslatorInterface $translator
     */
    public function __construct(TranslatorInterface $translator)
    {
        $this->translator = $translator;
    }

    /**
     * @param \DateInterval $diff
     *
     * @return string
     */
    public function format(\DateInterval $diff): string
    {
        foreach (DateTimeDiffFormatter::UNITS as $attr => $unit) {
            $count = $diff->$attr;
            if (0 !== $count) {
                return $this->getDiffMessage($count, (bool) $diff->invert, $unit);
            }
        }

        return $this->getNowMessage();
    }

    /**
     * @param \DateTimeInterface $from
     * @param \DateTimeInterface $to
     *
     * @return string
     */
    public function formatDiff(\DateTimeInterface $from, \DateTimeInterface $to): string
    {
        return $this->format($to->diff($from));
    }

    /**
     * @param int    $count
     * @param bool   $invert
     * @param string $unit
     *
     * @return string
     */
    private function getDiffMessage(int $count, bool $invert, string $unit): string
    {
        $id = sprintf('diff.%s.%s', $invert ? 'ago' : 'in', $unit);

        return $this->translator->trans($id, ['%count%' => $count], 'time');
    }

    /**
     * @return string
     */
    private function getNowMessage(): string
    {
        return $this->translator->trans('diff.now', [], 'time');
    }
}
