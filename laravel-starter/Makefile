SHELL := /bin/bash

.PHONY: build up stop down

build:
	docker compose -f docker-compose.yml build

up:
	docker compose -f docker-compose.yml up -d

stop:
	docker compose stop

down:
	docker compose down

migrate-fresh:
	docker compose exec -it app php artisan migrate:fresh --seed