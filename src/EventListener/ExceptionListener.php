<?php
/**
 * (c) Taranov Egor <dev@taranovegor.com>
 */

namespace App\EventListener;

use App\Exception\CatchableExceptionInterface;
use Monolog\Logger;
use Psr\Log\LoggerInterface;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\Event\ExceptionEvent;
use Symfony\Component\HttpKernel\Exception\BadRequestHttpException;
use Symfony\Component\HttpKernel\Exception\HttpException;
use Symfony\Contracts\Translation\TranslatorInterface;
use TelegramBot\Api\BotApi;
use TelegramBot\Api\InvalidJsonException;
use TelegramBot\Api\Types\Update;

/**
 * Class ExceptionListener
 */
class ExceptionListener
{
    /**
     * @var BotApi
     */
    private BotApi $api;

    /**
     * @var LoggerInterface
     */
    private LoggerInterface $logger;

    /**
     * @var TranslatorInterface
     */
    private TranslatorInterface $translator;

    /**
     * ExceptionListener constructor.
     *
     * @param BotApi              $api
     * @param LoggerInterface     $exceptionLogger
     * @param TranslatorInterface $translator
     */
    public function __construct(BotApi $api, LoggerInterface $exceptionLogger, TranslatorInterface $translator)
    {
        $this->api = $api;
        $this->logger = $exceptionLogger;
        $this->translator = $translator;
    }

    /**
     * @param ExceptionEvent $event
     *
     * @return void
     */
    public function onKernelException(ExceptionEvent $event)
    {
        $e = $event->getThrowable();
        $request = $event->getRequest();

        $responseCode = Response::HTTP_BAD_REQUEST;
        $logLevel = Logger::CRITICAL;
        $logContext = [
            'message' => $e->getMessage(),
            'response_code' => &$responseCode,
            'exception_class' => get_class($e),
            'exception_trace' => $e->getTraceAsString(),
            'during_processing' => false,
        ];

        try {
            if ($content = $request->getContent()) {
                $update = Update::fromResponse(
                    BotApi::jsonValidate($content, true)
                );
            } else {
                throw new BadRequestHttpException('Request data is empty');
            }

            if ($e instanceof CatchableExceptionInterface) {
                $logLevel = Logger::ERROR;

                if ($e->isProcessed()) {
                    $responseCode = Response::HTTP_OK;
                    $this->api->sendMessage(
                        $update->getMessage()->getChat()->getId(),
                        $this->translator->trans(
                            $e->getMessage(),
                            [],
                            'errors',
                            $update->getMessage()->getFrom()->getLanguageCode()
                        )
                    );
                }
            } else {
                $logLevel = Logger::CRITICAL;

                if ($e instanceof \TelegramBot\Api\HttpException) {
                    $logContext['parameters'] = [
                        $e->getParameters(),
                    ];
                }
            }
        } catch (\Throwable $exception) {
            $responseCode = Response::HTTP_OK;
            $logContext['during_processing'] = [
                'exception_class' => get_class($exception),
                'exception_message' => $exception->getMessage(),
            ];
        } finally {
            $this->logger->log(
                $logLevel,
                'on_kernel_exception: {response_code} {exception_class}[{message}]',
                $logContext
            );
            $event->setThrowable(new HttpException($responseCode));
            $event->setResponse(Response::create('', $responseCode));
        }
    }
}
