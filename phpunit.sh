#!/bin/sh

docker-compose exec php php bin/phpunit "$@"
