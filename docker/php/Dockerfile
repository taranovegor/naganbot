FROM php:7.4-fpm-alpine

ARG TZ=UTC
ARG UID=1000
ARG GID=1000

RUN apk update && apk add --no-cache shadow git libzip-dev libxml2-dev icu-dev

RUN docker-php-ext-install opcache json zip pdo_mysql bcmath intl

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

RUN usermod -u $UID www-data && groupmod -g $GID www-data

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

USER $UID:$GID

EXPOSE 9000
