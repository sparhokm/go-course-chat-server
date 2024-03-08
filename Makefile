RUN:=-f docker-compose.yml -f docker-compose-run.yml
DEBUG:=$(RUN) -f docker-compose-debug.yml
CLI:=docker-compose run --rm --no-deps cli

init: docker-down local-env-build \
	init-up wait-db db-migrations-up vendor-refresh

rebuild: docker-down local-env-build docker-build vendor-refresh

down: docker-down

run: init grpc-run
run-restart: down grpc-run

debug: init grpc-debug
debug-restart: down grpc-debug


init-up:
	docker-compose up -d

grpc-run:
	docker-compose $(RUN) up -d

grpc-debug:
	docker-compose $(DEBUG) up -d

local-env-build:
	chmod 777 ./docker/common/env-init.sh
	./docker/common/env-init.sh ./.env ./docker/.env ./docker/.local.env ./.env.local

docker-down:
	docker-compose down --remove-orphans

docker-down-clear:
	docker-compose down -v --remove-orphans

docker-pull:
	docker-compose pull

docker-build:
	docker-compose build --pull

wait-db:
	$(CLI) wait-for-it db:5432 -t 30

db-migrations-create:
	$(CLI) goose -dir migrations create $(filter-out $@,$(MAKECMDGOALS)) sql

db-migrations-status:
	$(CLI) goose -dir migrations status

db-migrations-up:
	$(CLI) goose -dir migrations up -v

db-migrations-down:
	$(CLI) goose -dir migrations down -v

vendor-refresh:
	$(CLI) go mod vendor

lint:
	$(CLI) golangci-lint run ./... --config .golangci.pipeline.yaml
	
generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/chat_v1
	$(CLI) protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=/go/bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=/go/bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto