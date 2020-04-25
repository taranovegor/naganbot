# Nagan bot

Create file for local environment
```
touch .env.local
```

Deploy application
```
sh redeploy.sh -p
```

Set telegram webhook
```
sh console.sh telegram:webhook:set <url> [<path-to-certificate>]
```

Stop all containers
```
sh redeploy.sh -s
```

## Manipulating application

Symfony console commands
```
sh console.sh <specific command>
```

Composer commands
```
sh composer.sh <specific command>
```

## Code quality

Check Symfony coding standards
```
sh checkstyle.sh
```
