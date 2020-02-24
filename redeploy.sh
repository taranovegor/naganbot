#!/bin/sh

docker-compose -f docker-compose.yml -f docker-compose.dev.yml stop
docker-compose -f docker-compose.yml -f docker-compose.dev.yml build
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

docker-compose exec php composer install
docker-compose exec php php /var/www/html/bin/console rabbitmq-supervisor:rebuild
sleep 5
docker-compose exec php php /var/www/html/bin/console rabbitmq-supervisor:control --wait-for-supervisord start
docker-compose exec php php bin/console doctrine:migrations:migrate --allow-no-migration --no-interaction
