#!/bin/sh

if [ -f .env.local ]; then
    export $(cat .env.local)
fi

STOP=false
PROD=false
DOCKER_COMPOSE_FILES="-f docker-compose.yml"

while getopts "sp" OPTION; do
    case $OPTION in
        s)
            STOP=true
            ;;
        p)
            PROD=true
            ;;
    esac
done

if [ $PROD = false ]; then
    DOCKER_COMPOSE_FILES="$DOCKER_COMPOSE_FILES -f docker-compose.dev.yml"
fi

docker-compose $DOCKER_COMPOSE_FILES stop
if [ $STOP = true ]; then
    exit
fi
docker-compose $DOCKER_COMPOSE_FILES build
docker-compose $DOCKER_COMPOSE_FILES up -d

docker-compose exec php composer install
docker-compose exec php php bin/console doctrine:migrations:migrate --allow-no-migration --no-interaction
