#!/usr/bin/env bash
set -eoux pipefail

function runLinters() {
  go install honnef.co/go/tools/cmd/staticcheck@latest
  go fmt ./...
  go vet ./...
  staticcheck -tags=withDocker ./...
}

function waitForPostgres() {
  until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -c '\q' -p "$POSTGRES_PORT"; do
      >&2 echo "postgres is unavailable - sleeping ..."
      sleep 1
  done
}

function runTests() {
  CGO_ENABLED=1 go test -v -race -cover -p=1 ./... -mod=vendor -tags=withDocker
}

function main() {
  runLinters
  set +euox pipefail
  waitForPostgres
  set -euox pipefail
  runTests
}

main "$@"