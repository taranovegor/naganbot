#!/bin/sh

docker-compose exec php php /var/www/html/bin/console $1 $2 $3 $4
