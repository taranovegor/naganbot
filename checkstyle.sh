#!/bin/sh

docker run --rm -v $(pwd):/data cytopia/phpcs:3 --standard=phpcs.xml
docker run --rm -v $(pwd):/app phpstan/phpstan analyse -l 1 -c /app/phpstan.neon /app/src
