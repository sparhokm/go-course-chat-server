PROJECT = github.com/sparhokm/go-course-ms-chat-server

RUN:=-f docker-compose.yml -f docker-compose-run.yml
DEBUG:=$(RUN) -f docker-compose-debug.yml
CLI:=docker-compose run --rm --no-deps cli

init: local-env-build docker-down \
	init-up wait-db db-migrations-up vendor-refresh

rebuild: local-env-build docker-down docker-pull docker-build vendor-refresh

down: docker-down

run: grpc-run
run-restart: down grpc-run

debug: grpc-debug
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
	$(CLI) go mod tidy;
	$(CLI) go mod vendor;
	make vendor-proto

vendor-proto:
	@if [ ! -d vendor.protogen/validate ]; then \
		mkdir -p vendor.protogen/validate &&\
		$(CLI) git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
		mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
		rm -rf vendor.protogen/protoc-gen-validate ;\
	fi
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p  vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		$(CLI) git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi

fmt:
	$(CLI) go fmt ./...

lint:
	$(CLI) golangci-lint run ./... --config .golangci.pipeline.yaml
	
generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/chat_v1
	mkdir -p pkg/swagger
	$(CLI) protoc --proto_path api/chat_v1 --proto_path vendor.protogen \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=/go/bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=/go/bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/chat_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=/go/bin/protoc-gen-validate \
	--grpc-gateway_out=pkg/chat_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=/go/bin/protoc-gen-grpc-gateway \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=/go/bin/protoc-gen-openapiv2 \
	api/chat_v1/chat.proto

mockery:
	$(CLI) mockery

test-coverage:
	$(CLI) go test ./... -coverprofile=coverage.tmp.out -covermode=count -coverpkg=$(PROJECT)/... -count=5
	grep -v 'mocks\|config\|/pkg/chat_v1\|/cmd\|/internal/app' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	$(CLI) go tool cover -func=./coverage.out | grep "total";
	rm coverage.out

test-coverage-ci:
	go test ./... -coverprofile=coverage.tmp.out -covermode=atomic -coverpkg=$(PROJECT)/... -race -count=5
	grep -v 'mocks\|config\|/pkg/chat_v1\|/cmd\|/internal/app' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out