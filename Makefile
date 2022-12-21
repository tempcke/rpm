project = $(shell basename $(shell pwd))

help:				## display help information
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build:				## create binary
	go build -o bin/rpm cmd/rpmserver/*.go

run: .env build dockerUp	## create and run binary
	godotenv bin/rpm

check: test lint	## lint + test, pre-commit hook

lint:				## fmt, vet, and staticcheck
	test -z $(shell go fmt ./...)
	go vet ./...
	staticcheck -tags=withDocker ./...

test:				## execute tests
	godotenv time -p go test -failfast -race -count=1 ./... -cover | grep -v '\[no test'

testAll: dockerUp	## run all tests including those that need docker/postgres
	godotenv time -p go test -failfast -p=1 -count=1 ./... -tags=withDocker -cover | grep -v '\[no test'

testCI:				## exact tests the way buildkite does, use for local debug of buildkite failure
	docker-compose -f docker-compose-ci.yml -p $(project)-ci run --rm appci /bin/sh -e -c 'bash pipeline/test.sh' || true
	docker-compose -f docker-compose-ci.yml -p $(project)-ci down

dockerUp: init		## docker-compose up
	@if [ ! "$(shell docker-compose ps --services --filter "status=running")" = "postgres" ]; then \
		docker-compose up -d; \
		sleep 3; \
	fi

dockerDown:	init	## docker-compose down
	docker-compose down

dockerRestart: dockerDown dockerUp	## dockerDown && dockerUp

clean: dockerDown	## dockerDown && docker-compose down for CI
	docker-compose -f docker-compose-ci.yml -p $(project)-ci down

init: .env .git/hooks/pre-commit

.env: .git/hooks/pre-commit ## copy .env.example to .env
	cp .env.example .env

.git/hooks/pre-commit:
	cp -r .githooks/* .git/hooks/

.PHONY: help check lint test testAll testCI dockerUp dockerDown dockerRestart clean init
