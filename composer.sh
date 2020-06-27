#!/bin/sh

docker-compose exec php php /usr/bin/composer "$@"
