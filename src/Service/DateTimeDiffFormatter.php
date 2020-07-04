<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\Service;

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
