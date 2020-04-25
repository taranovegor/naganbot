#!/bin/sh

docker-compose exec php php -d memory_limit=2g /usr/bin/composer $1 $2 $3 $4
