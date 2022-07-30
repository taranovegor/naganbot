# naganbot
Shoot yourself in Russian roulette

## Set up
### Requirements
- Go 1.20 or higher

### Configuration
Copy an instance of the environment file and save it as a file `.env`
```shell
cp .env.dist .env
```
... and configure as you need. All configuration instructions are given in the configuration example

### Launch

Install the dependencies in the vendor directory and run the application
```shell
go mod vendor
go build -v ./...
go run -v ./...
```
Enjoy!
