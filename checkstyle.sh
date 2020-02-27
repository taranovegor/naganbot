#!/bin/sh

docker-compose exec php php vendor/bin/phpcs --standard=vendor/escapestudios/symfony2-coding-standard/Symfony --ignore=src/Kernel.php,src/Migrations/* -q src
docker run --rm -v $(pwd):/app phpstan/phpstan analyse -l 1 -c /app/phpstan.neon /app/src
