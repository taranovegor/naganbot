FROM php:7.4-fpm

RUN apt-get update && \
    apt-get -y --no-install-recommends install \
        git \
        libzip-dev \
        supervisor \
        libxml2-dev
RUN docker-php-ext-install \
        opcache \
        json \
        zip \
        pdo_mysql \
        bcmath

COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

EXPOSE 9000
