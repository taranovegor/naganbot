<?php

namespace App;

use Dotenv\Dotenv;
use unreal4u\TelegramAPI\HttpClientRequestHandler;
use unreal4u\TelegramAPI\Telegram\Methods\SendMessage;
use unreal4u\TelegramAPI\Telegram\Types\Update;
use unreal4u\TelegramAPI\TgLog;

class Kernel
{
    /**
     * @var string
     */
    private $request;
    
    /**
     * @var \unreal4u\TelegramAPI\Telegram\Types\Update;
     */
    private $update;
    
    /**
     * @var \unreal4u\TelegramAPI\TgLog
     */
    private $bot;
    
    /**
     * @throws \Exception
     */
    public function register()
    {
        (Dotenv::create(realpath(__DIR__ . '/../')))->load();
        
        
        
        
        
        
        
        
        
        $loop = \React\EventLoop\Factory::create();
        $bot = new TgLog($_ENV['BOT_TOKEN'], new HttpClientRequestHandler($loop));
    
        $updateData = json_decode(file_get_contents('php://input'), true);
    
        $update = new Update($updateData);
    
        $sendMessage = new SendMessage();
        $sendMessage->chat_id = $update->message->chat->id;
        $sendMessage->text = 'Hello world!';
    
        $bot->performApiRequest($sendMessage);
        $loop->run();
    }
    
    /**
     * @throws \Exception
     */
    private function route()
    {
//        $router = new Router();
        
//        $this->update->getMessage()
        
//        $action = new SendMessage(
//            $this->update->getMessage()->getChat()->getId(),
//            ".". ($this->update->getMessage()-> ? 1 : 0)
//        );
//
//        $method = end(explode('\\', get_class($action)));
//        if (method_exists($this->bot, $method)) {
//            $this->bot->{$method}($action);
//        } else {
//            // log it
//        }
    }
    
    /**
     * Run instance of app
     */
    public function run()
    {
        try {
//            $this->request = file_get_contents('php://input');
//            $this->update = Update::create(json_decode($this->request, true));
            $this->register();
            $this->route();
        } catch (\Exception $e) {
            echo $e->getMessage();
        }
    }
}
