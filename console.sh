#!/bin/sh

docker-compose exec php php /var/www/html/bin/console "$@"
