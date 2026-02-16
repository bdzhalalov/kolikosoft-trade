build:
	docker-compose build

run:
	docker-compose up -d
	sleep 1
	docker-compose ps

start: build run
	@echo "All services are up"

stop:
	docker-compose stop