# naganbot
Shoot yourself in Russian roulette

## Set up
### Requirements
- Docker, docker-compose
- MySQL 8/MariaDB 11

### Configuration
Copy an instance of the environment file and save it as a file `.env`
```shell
cp .env.dist .env
```
... and configure as you need. Environment variables are described in comments

Build or pull docker images of application
```shell
make container-build
```
```shell
make container-pull
```

### Launch
Run the application using the Make tool
```shell
make start
```
