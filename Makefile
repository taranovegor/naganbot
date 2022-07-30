.PHONY: help

DOCKER_COMPOSE_OPTIONS = -f docker-compose.yaml
DOCKER_COMPOSE = docker-compose $(DOCKER_COMPOSE_OPTIONS)

help: ## Displays help for a command
	@printf "\033[33mUsage:\033[0m\n  make [options] [target] ...\n\n\033[33mAvailable targets:%-13s\033[0m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' 'Makefile' | awk 'BEGIN {FS = ":.*?## "}; {printf "%-2s\033[32m%-17s\033[0m %s\n", "", $$1, $$2}'

container-build: ## Builds the application's docker containers
	$(DOCKER_COMPOSE) build --compress --force-rm

container-up: ## Launches docker application containers
	$(DOCKER_COMPOSE) up --detach --remove-orphans --force-recreate
	$(DOCKER_COMPOSE) ps

run: ## executes the application launch
	$(MAKE) container-build
	$(MAKE) container-up
