.PHONY: all build-web run clean dev-db

DATABASE_URL ?= postgres://user:password@localhost:5432/opsdeck?sslmode=disable

all: build-web run

DOCKER_COMPOSE ?= docker compose

dev-db:
	$(DOCKER_COMPOSE) up -d db redis

build-web:
	cd web && npm install && npm run build
	mkdir -p cmd/server/dist
	cp -r web/dist/* cmd/server/dist/

run: build-web dev-db
	@echo "Waiting for DB..."
	@sleep 2
	DATABASE_URL=$(DATABASE_URL) go run cmd/server/main.go

clean:
	rm -rf web/dist
	rm -rf cmd/server/dist
	rm -f server.log server.pid
