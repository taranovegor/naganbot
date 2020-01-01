<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\MessageBuilder;

use Symfony\Contracts\Translation\TranslatorInterface;

/**
 * Class AbstractTranslatableMessageBuilder
 */
abstract class AbstractTranslatableMessageBuilder extends AbstractMessageBuilder
{
    /**
     * @var TranslatorInterface
     */
    protected TranslatorInterface $translator;

    /**
     * AbstractTranslatableMessageBuilder constructor.
     *
     * @param TranslatorInterface $translator
     */
    public function __construct(TranslatorInterface $translator)
    {
        parent::__construct();
        $this->translator = $translator;
    }
}
