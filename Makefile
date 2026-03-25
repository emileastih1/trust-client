.PHONY: server proto lint clean .FORCE
COVERAGE_DIR=coverage
GO_TARGETS=./internal/... ./cmd/...
GOFLAGS=

.PHONY: test

all:
	@echo BBAPI_STAGE is $$BBAPI_STAGE
	@echo BBAPI_ADDRESS_SECRET is $$BBAPI_ADDRESS_SECRET
	@echo BBAPI_DB_URL is $$BBAPI_DB_URL

deps:
	go get ./...
	go mod tidy

test:
	@go test ${GO_TARGETS}

load-test:
	@PYTHONPATH=test/integration_test python3 test/integration_test/use_case_test.py --number_of_threads 80 --total_number_of_iterations 100

load-test-ci:
	@PYTHONPATH=test/integration_test python3 test/integration_test/use_case_test.py --number_of_threads 20 --total_number_of_iterations 100


integration-test:
	@PYTHONPATH=test/integration_test python3 -m pytest

cover:
	mkdir -p ${COVERAGE_DIR}
	go test -coverprofile=${COVERAGE_DIR}/coverage.out ${GO_TARGETS} && \
		go tool cover -html=${COVERAGE_DIR}/coverage.out -o ${COVERAGE_DIR}/coverage.html

lint:
	golangci-lint run
	buf lint

format-test:
	autopep8 --in-place --recursive test/integration_test/*_test.py

generate:
	buf generate
	wire ./...

vet:
	go vet ./...

compile: deps generate lint vet
	go build -o ./build/server ${GOFLAGS} ./main.go

build:
	go build -o ./build/server ${GOFLAGS} ./main.go

init-secrets: export BBAPI_STAGE=development
init-secrets: export BBAPI_ADDRESS_SECRET=D2R5WSpZZyj0RSw1F7HrLVDxZX3nrX4Xb6+ce4qX6H21y8bRqtq/S3SNYdL5NJ0bHY/uPO5DOjl0XJaNRKXRsg==
init-secrets: export BBAPI_DB_URL=postgres://postgres:password@127.0.0.1:5432/my_postgres_db?sslmode=disable&connect_timeout=5
init-secrets: all
	go run cmd/bootstrap/main.go

local-dep: export BBAPI_STAGE=local
local-dep: export BBAPI_ADDRESS_SECRET=D2R5WSpZZyj0RSw1F7HrLVDxZX3nrX4Xb6+ce4qX6H21y8bRqtq/S3SNYdL5NJ0bHY/uPO5DOjl0XJaNRKXRsg==
local-dep: export BBAPI_DB_URL=postgres://postgres:password@127.0.0.1:5432/my_postgres_db?sslmode=disable&connect_timeout=5
local-dep: all
	docker compose up -d
	go run cmd/db/main.go up

server: export BBAPI_STAGE=local
server: export BBAPI_ADDRESS_SECRET=D2R5WSpZZyj0RSw1F7HrLVDxZX3nrX4Xb6+ce4qX6H21y8bRqtq/S3SNYdL5NJ0bHY/uPO5DOjl0XJaNRKXRsg==
server: export BBAPI_DB_URL=postgres://postgres:password@127.0.0.1:5432/my_postgres_db?sslmode=disable&connect_timeout=5
server: all
	@go build -o ./build/server ${GOFLAGS} ./main.go && ./build/server

clean:
	go run cmd/db/main.go reset
	docker compose down --volumes --remove-orphans

docker-build:
	docker build -t local-bbapi .

docker-up:
	docker compose -f docker-compose.local.yaml up -d

docker-all:
	docker compose -f docker-compose-all.yaml up

.FORCE:
