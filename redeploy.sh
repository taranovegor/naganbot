#!/bin/sh

if [ -f "/sys/class/dmi/id/product_name" ]; then
  SYS_NAME=$(cat /sys/class/dmi/id/product_name)
else
  SYS_NAME="unknown"
fi

if [ $(command -v vagrant) ]; then
  VAGRANT_SUPPORT=true
else
  VAGRANT_SUPPORT=false
fi

STOP=false
DOWN=false
KILL=false
PROD=false
VAGRANT=false
DOCKER_COMPOSE_FILES="-f docker-compose.yaml"

for ARG in "$@"; do
  case $ARG in
    -s)
      STOP=true;
      ;;
  esac
done
exit;

if [ $VAGRANT = true ]; then
    vagrant ssh -c "cd /home/vagrant/workspace && sh redeploy.sh $@"
    exit
fi

if [ -f .env.local ]; then
    export $(cat .env.local | sed "s/#.*//g")
fi

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
