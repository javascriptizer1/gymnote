include .env

MIGRATIONS_SRC_DIR := migrations/clickhouse

docker-up:
	docker-compose up -d 

run:
	@set -a; \
	source .env; \
	go run ./cmd/gymnote/main.go

migrate-up:
	GOOSE_DRIVER=clickhouse \
	GOOSE_DBSTRING="tcp://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}" \
	GOOSE_MIGRATION_DIR=${MIGRATIONS_SRC_DIR} \
	goose up

# make migrate-gen name=<name>
migrate-gen:
	goose -dir $(MIGRATIONS_SRC_DIR) create ${name} sql

fmt:
	gofmt -w .
	goimports -w -local github.com/javascriptizer1 .