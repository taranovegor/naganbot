<?php

namespace App;

class Router
{
    /**
     * @var string
     */
    private $route;
    
    /**
     * @var array
     */
    private $routes;
    
    /**
     * @return array
     */
    public function getRoutes(): array
    {
        if (empty($this->routes)) {
            $this->routes = require_once '../router/routes.php';
        }
        return $this->routes;
    }
    
    /**
     * @param string $command
     * @throws \Exception
     */
    public function search(string $command)
    {
        $command = ltrim($command, '/');
    
        if (array_key_exists($command, $this->getRoutes()) !== true) {
            throw new \Exception('route for command '. $command .' not found');
        }
    
        $this->route = $this->getRoutes()[$command];
    }
    
    /**
     *
     */
    public function dispatch()
    {
        new $this->route();
        
    }
    
    /**
     * @param string $command
     * @return \App\Router
     * @throws \Exception
     */
    public static function create(string $command): self
    {
        $router = new self;
        $router->search($command);
        return $router;
    }
}