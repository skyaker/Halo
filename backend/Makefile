include .env
export

SHELL := /bin/bash

PROJECT_NAME := habitlogger
DATABASES := docker-compose-databases.yml
BACKEND := docker-compose-backend.yml
KAFKA := docker-compose-kafka.yml

COMPOSE := docker compose -p $(PROJECT_NAME) -f $(DATABASES)
COMPOSE_BACKEND := docker compose -p $(PROJECT_NAME) -f $(BACKEND)
COMPOSE_KAFKA := docker compose -p $(PROJECT_NAME) -f $(KAFKA)

# DIRS = note_service category_service auth_service
DIRS = note_service user_service auth_service

.PHONY: all
all: down-all network-all build-all up-all

.PHONY: down-all
down-all: down-backend down-databases down-kafka

.PHONY: build-all
build-all: build-databases build-kafka build-backend

.PHONY: up-all
up-all: up-databases up-kafka up-backend

.PHONY: clean-all
clean-all:
	@docker container prune -f

.PHONY: network-all
network-all: rm-network create-network

.PHONY: create-network
create-network:
	@docker network create backend
	@docker network create frontend

.PHONY: rm-network
rm-network:
	@docker network rm backend || true
	@docker network rm frontend || true

.PHONY: hard-rm-network
hard-rm-network:
	@docker rm -f $$(docker ps -aq) || true
	@docker volume prune -f || true
	@docker network prune -f || true
	@docker network rm backend frontend || true

# -------------------------------------------KAFKA-------------------------------------------

.PHONY: up-kafka
up-kafka:
	$(COMPOSE_KAFKA) up -d --build
	@make kafka-topic-init

.PHONY: build-kafka
build-kafka:
	$(COMPOSE_KAFKA) build --no-cache --force-rm

.PHONY: down-kafka
down-kafka:
	$(COMPOSE_KAFKA) down || true

.PHONY: kafka-topic-init
kafka-topic-init:
	@docker exec kafka ./opt/kafka/bin/kafka-topics.sh --create --topic user-created --bootstrap-server kafka:9092 --partitions 1 --replication-factor 1 --if-not-exists

# -------------------------------------------DATABASE-------------------------------------------

.PHONY: clean-postgres
clean-postgres:
	@for dir in $(DIRS); do \
		if [ "$$dir" = "note_service" ]; then \
			DB_NAME=$(NOTE_DB); \
			CONTAINER_NAME=$(NOTE_POSTGRES); \
		fi; \
		docker rm -f $$CONTAINER_NAME || true; \
	done

.PHONE: clean-auth-redis
clean-auth-redis:
	@docker exec -it $(AUTH_REDIS) redis-cli flushall

.PHONY: down-databases
down-databases:
	$(COMPOSE) down || true

.PHONY: build-databases
build-databases:
	$(COMPOSE) build --no-cache --force-rm

.PHONY: up-databases
up-databases:
	$(COMPOSE) up -d --build

.PHONY: pg-seed-data
pg-seed-data:
	@for dir in $(DIRS); do \
		if [ "$$dir" = "note_service" ]; then \
			DB_NAME=$(NOTE_DB); \
			CONTAINER_NAME=$(NOTE_POSTGRES); \
		fi; \
		docker exec -i $$CONTAINER_NAME psql -U $(POSTGRES_USER) -d $$DB_NAME < $$dir/db/init/seed_data.sql; \
	done

# -------------------------------------------SERVICES-------------------------------------------

.PHONY: down-backend
down-backend:
	$(COMPOSE_BACKEND) down || true

.PHONY: build-backend
build-backend:
	$(COMPOSE_BACKEND) build

.PHONY: up-backend
up-backend:
	$(COMPOSE_BACKEND) up --build -d

# -------------------------------------------FILES------------------------------------------- 

.PHONY: update-env
update-env:
	@for dir in $(DIRS); do \
	    rm -f $$dir/.env; \
	    cp .env $$dir/.env; \
	done

# -------------------------------------------KAFKA------------------------------------------- 

.PHONY: reset-kafka
reset-kafka:
	@echo "Stopping containers..."
	docker-compose stop kafka
	@echo "Removing Kafka logs inside container..."
	docker exec kafka rm -rf /tmp/kraft-combined-logs/*
	@echo "Starting Kafka again..."
	docker-compose start kafka

.PHONY: wait
wait:
	./wait.sh
