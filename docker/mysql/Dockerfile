FROM mariadb:10.5 AS build

ARG TZ=UTC
ARG UID=1000
ARG GID=1000

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN usermod -u $UID mysql && groupmod --non-unique -g $GID mysql
