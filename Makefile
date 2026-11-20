include .env
export

COMPOSE = docker compose --env-file .env -f deploy/docker-compose.yml

.PHONY: up down migrate status run

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

migrate:
	goose -dir migrations postgres "$(DATABASE_URL)" up

status:
	goose -dir migrations postgres "$(DATABASE_URL)" status

run:
	go run ./cmd/coordinator