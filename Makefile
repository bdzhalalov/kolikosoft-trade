-include .env
export
MIGRATIONS_DIR = /var/www/migrations
PG_DSN = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
TEST_CONTAINER_NAME := kolikosoft-trade-test
.PHONY: test build run migrate-up migrate-down migrate-create

build:
	docker-compose build

run:
	docker-compose up -d
	sleep 1
	docker-compose ps

start: build run migrate-up
	@echo "All services are up"

migrate-up:
	docker exec -it kolikosoft-trade migrate \
      	-path $(MIGRATIONS_DIR) \
    	-database "$(PG_DSN)" up

migrate-down:
	docker exec -t kolikosoft-trade migrate \
        -path $(MIGRATIONS_DIR) \
        -database "$(PG_DSN)" down 1

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate create name=add_table_x"; exit 1)
	docker exec -it kolikosoft-trade migrate create \
	-ext=sql -seq -dir=$(MIGRATIONS_DIR) $(name)

test:
	docker build --target builder -t $(TEST_CONTAINER_NAME) .
	docker run --rm $(TEST_CONTAINER_NAME) go test -v ./...

stop:
	docker-compose stop