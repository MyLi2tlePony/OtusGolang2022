BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"

MAIN_CALENDAR := "./cmd/calendar"
MAIN_SCHEDULER := "./cmd/scheduler"
MAIN_SENDER := "./cmd/sender"

CONFIG_CALENDAR := "./configs/calendar/config.toml"
CONFIG_SCHEDULER := "./configs/scheduler/config.toml"
CONFIG_SENDER := "./configs/sender/config.toml"

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc --go_out=internal/server/grpc --go-grpc_out=internal/server/grpc api/CalendarService.proto

rabbitmq:
	docker run --name rabbitmq --detach --rm -p 15672:15672 -p 5672:5672 rabbitmq:3-management

postgres:
	docker run --name postgresdb --env POSTGRES_PASSWORD="1234512345" --publish "5436:5432" --detach --rm postgres

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" $(MAIN_CALENDAR)
	go build -o $(BIN_SENDER) $(MAIN_SENDER)
	go build -o $(BIN_SCHEDULER) $(MAIN_SCHEDULER)

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" $(MAIN_CALENDAR)

build-scheduler:
	go build -o $(BIN_SCHEDULER) $(MAIN_SCHEDULER)

build-sender:
	go build -o $(BIN_SENDER) $(MAIN_SENDER)

run-calendar: build-calendar
	$(BIN_CALENDAR) -config $(CONFIG_CALENDAR)

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) -config $(CONFIG_SCHEDULER)

run-sender: build-sender
	$(BIN_SENDER) -config $(CONFIG_SENDER)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build-calendar
	$(BIN_CALENDAR) version

migrate-up:
	docker exec -it postgresdb createdb --username postgres calendardb
	migrate -path migrations -database "postgresql://postgres:1234512345@localhost:5436/calendardb?sslmode=disable" -verbose up

.PHONY: generate rabbitmq postgre\
		install-lint-deps lint test \
		build build-calendar build-scheduler build-sender \
 		run-calendar run-scheduler run-sender \
 		build-img run-img version migrateup
