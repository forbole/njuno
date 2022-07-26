VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

export GO111MODULE = on

###############################################################################
###                                   All                                   ###
###############################################################################

all: lint test-unit install

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/forbole/njuno/cmd.Version=$(VERSION) \
	-X github.com/forbole/njuno/cmd.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	@echo "building njuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/njuno.exe ./cmd/njuno
else
	@echo "building njuno binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/njuno ./cmd/njuno
endif
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "installing njuno binary..."
	@go install -mod=readonly $(BUILD_FLAGS) ./cmd/njuno && cp modules/staking/utils/validators_query.sh ~/.njuno
.PHONY: install

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

stop-docker-test:
	@echo "Stopping Docker container..."
	@docker stop njuno-test-db || true && docker rm njuno-test-db || true
.PHONY: stop-docker-test

start-docker-test: stop-docker-test
	@echo "Starting Docker container..."
	@docker run --name njuno-test-db -e POSTGRES_USER=njuno -e POSTGRES_PASSWORD=password -e POSTGRES_DB=njuno -d -p 6433:5432 postgres
.PHONY: start-docker-test

coverage:
	@echo "viewing test coverage..."
	@go tool cover --html=coverage.out
.PHONY: coverage

test-unit: start-docker-test
	@echo "Executing unit tests..."
	@go test -mod=readonly -v -coverprofile coverage.txt ./...
.PHONY: test-unit

lint:
	golangci-lint run --out-format=tab
.PHONY: lint

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint-fix

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -name '*.pb.go' | xargs goimports -w -local github.com/forbole/njuno
.PHONY: format

clean:
	rm -f tools-stamp ./build/**
.PHONY: clean
