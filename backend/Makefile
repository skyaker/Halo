SHELL := /bin/bash

PROJECT_NAME := habitlogger
DATABASES := docker-compose-databases.yml
# BACKEND := docker-compose-backend.yml

COMPOSE := docker compose -p $(PROJECT_NAME) -f $(DATABASES)

.PHONY: all
all: down-all build-all up-all

.PHONY: down-all
down-all: down-databases

.PHONY: build-all
build-all: build-databases

.PHONY: up-all
up-all: up-databases

.PHONY: clean-all
clean-all:
	docker container prune -f

# -------------------------------------------DATABASE-------------------------------------------

.PHONY: clean-postgres
clean-postgres:
	docker rm -f habitlogger-postgres || true

.PHONY: down-databases
down-databases:
	$(COMPOSE) down || true

.PHONY: build-databases
build-databases:
	$(COMPOSE) build --no-cache --force-rm

.PHONY: up-databases
up-databases:
	$(COMPOSE) up -d --build
