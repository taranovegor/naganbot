#!/bin/sh

if [ -f .env.local ]; then
    export $(cat .env.local | sed "s/#.*//g")
fi

STOP=false
DOWN=false
KILL=false
PROD=false
DOCKER_COMPOSE_FILES="-f docker-compose.yaml"

while getopts ":sdkp" OPTION; do
    case $OPTION in
        s)
            STOP=true
            ;;
        d)
            DOWN=true
            ;;
        k)
            KILL=true
            ;;
        p)
            PROD=true
            ;;
    esac
done

if [ $PROD = false ]; then
    DOCKER_COMPOSE_FILES="$DOCKER_COMPOSE_FILES -f docker-compose.dev.yaml"
fi

if [ $KILL = true ]; then
    docker-compose $DOCKER_COMPOSE_FILES kill
elif [ $DOWN = true ]; then
    docker-compose $DOCKER_COMPOSE_FILES down
else
    docker-compose $DOCKER_COMPOSE_FILES stop
fi

if [ $STOP = true ]; then
    exit
fi

docker-compose $DOCKER_COMPOSE_FILES build
docker-compose $DOCKER_COMPOSE_FILES up -d

docker-compose exec php composer install
docker-compose exec php php bin/console doctrine:migrations:migrate --allow-no-migration --no-interaction
