#!/bin/sh

docker-compose -f docker-compose.yml stop
docker-compose -f docker-compose.yml build
docker-compose -f docker-compose.yml up -d

docker exec -it tg_bot_russian_roulette_php composer install
docker exec -it tg_bot_russian_roulette_php php /var/www/html/bin/console rabbitmq-supervisor:rebuild
sleep 5
docker exec -it tg_bot_russian_roulette_php php /var/www/html/bin/console rabbitmq-supervisor:control --wait-for-supervisord start
docker exec -it tg_bot_russian_roulette_php php bin/console doctrine:migrations:migrate --allow-no-migration --no-interaction
