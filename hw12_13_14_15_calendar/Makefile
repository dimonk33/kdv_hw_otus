BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	cd migrations || goose postgres "user=user password=password  dbname=hw sslmode=disable" up || cd ..

generate:
	cd api || protoc -I .  --go_out=. --go-grpc_out=. --grpc-gateway_out ./gen --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true EventService.proto || cd ..

.PHONY: build run build-img run-img version test lint
