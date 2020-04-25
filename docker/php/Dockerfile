FROM php:7.4-fpm-alpine

RUN apk update && apk add --no-cache git libzip-dev libxml2-dev icu-dev

RUN docker-php-ext-install opcache json zip pdo_mysql bcmath intl

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

EXPOSE 9000
