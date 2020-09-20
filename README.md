# Nagan bot

## Application setup

Clone project
```
git clone git@github.com:taranovegor/nagan-bot
```

Create a local environment variables file
```
touch .env.local
```

To launch containers and launch the application, run the command
```
sh redeploy.sh
```

The above command `redeploy.sh` has arguments

| arg | description               
| --- | ---
| s   | stop application    
| d   | down application containers
| k   | kill application containers 
| p   | run only prod containers

Set webhook telegram bot
```
sh console.sh telegram:webhook:set <url> [<path-to-certificate>]
```

### Setup for development

For local development tools (like Xdebug) to work correctly, you need to get host address of machine using command below and specify result as value of `HOST_MACHINE` environment variable
```
hostname -I | head -n1 | cut -d " " -f1
```

#### Xdebug

In PHPStorm go to: `Languages & Frameworks` > `PHP` > `Servers` > and set the following settings:
![Server configuration](https://mattermost.branderstudio.com/files/ny5r5544ctfadkkareyj33p9uy/public?h=s1Xo6e44eJgyk4gBwcQmJQP9mVkFOjwECm-yUJ3RF4E)

To use a different server configuration, change the name of the environment variable `XDEBUG_SERVER_NAME`

## Application management

Symfony console commands
```
sh console.sh [command]
```

Composer commands
```
sh composer.sh [options] [command]
```

## Code quality

Check Symfony code style 
```
sh checkstyle.sh
```

Run PHPUnit
```
sh phpunit.sh
```
