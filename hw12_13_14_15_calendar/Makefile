BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar && go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler && go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN_CALENDAR) -config ./configs/config.toml &&

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_CALENDAR) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	cd migrations && goose postgres "user=user password=password  dbname=hw sslmode=disable" up ; cd ..

generate:
	cd api && mkdir gen ; go generate ; cd ..

.PHONY: build run build-img run-img version test lint migrate generate
