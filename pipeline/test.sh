#!/usr/bin/env bash
set -eoux pipefail

runLinters() {
# commented out for now because go install fails: lookup proxy.golang.org: Temporary failure in name resolution
#  go install honnef.co/go/tools/cmd/staticcheck@latest
#  staticcheck -tags=withDocker ./...
  go fmt ./...
  go vet ./...
}

waitForPostgres() {
#  this requires psql to be installed in the container and it currently isn't ...
#  until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -c '\q' -p "$POSTGRES_PORT"; do
#      >&2 echo "postgres is unavailable - sleeping ..."
#      sleep 1
#  done
  sleep 1
}

runTests() {
  go test -failfast -p=1 -count=1 ./... -tags=withDocker -cover | grep -v '\[no test'
}

main() {
  runLinters
  set +euox pipefail
  waitForPostgres
  set -euox pipefail
  runTests
}

main "$@"