#!/usr/bin/env bash

# Usage: make apiCheck
#
# DO NOT run this file directly because the ID's used are not dynamic and it will error on 2nd run
# why not make them dynamic you ask?  even if they were dynamic then the row count in the db would grow with each run
# meaning the output list would get longer.  The purpose of this file is to generate a clean easy to read apicheck.md file
# which will later be used as a reference when creating swagger.  So lets not over complicate it for now.
#
# Troubleshooting
# If this isn't working for you it may be that the port is in use by a process that isn't being killed
# try running: lsof -i :8080
# look for the PID of the process and then run: kill -9 <PID>

outputMD=apicheck.md
outputLog=apicheck.log
exitCode=0
outputErrorLog=apicheck-errors.log

h1() {
  {
    echo ""
    echo "# $*"
  } > ${outputMD}
}
h2() {
  p ""
  p "## $*"
}
h3() {
  p ""
  p "### $*"
}
p() {
  echo "$*" >> "${outputMD}"
  v "$*"
}
v() {
  if [ "$LOG_LEVEL" = "debug" ]; then
    echo "$*"
  fi
}

execCmd() {
  h3 "$1"
  echo "\`\`\`" >> ${outputMD}
  for var in "${@:2}"; do
    cmdOutput=$(sh -c "$var")
    cmdExitCode=$?
    p "$var"
    p "$cmdOutput"
    if [ $cmdExitCode -ne 0 ]; then
      p "ERROR!"
      echo "ERROR: $1 exited with $cmdExitCode" >> $outputErrorLog
      exitCode=$cmdExitCode
    fi
  done
  echo "\`\`\`" >> ${outputMD}
}
freePort() {
  lsof -i :8080 -sTCP:LISTEN | awk 'NR > 1 {print $2}' | xargs kill -15 > /dev/null 2>&1
}


main() {
  echo "starting server ..."
  freePort
  echo "" > $outputLog
  echo "" > $outputErrorLog
  nohup godotenv go run ./cmd/rpmserver/main.go > "$outputLog" & pid=$!
  sleep 5
  echo "started server PID: $pid"
  echo "cleaning up database entities"
  echo "starting curl requests ..."
  h1 "opp-inventory API check"
  p 'This is a generated document, to re-generate run: make apiCheck'

  execRoutes

  echo "shutting down server ..."
  echo ""
  freePort
  echo "Finished, to view results run: less ${outputMD}"
  if [ $exitCode -ne 0 ]; then
    cat $outputErrorLog
  fi
  echo ""

  exit $exitCode
}

execRoutes() {
  h2 "Property CRUD"
  execCmd "PUT property1" "${putProperty1} | json_pp"
  execCmd "GET property1" "${getProperty1} | json_pp"
  execCmd "PUT and DELETE property2" "${putProperty2}" "${delProperty2}"
  execCmd "GET properties" "${getProperties} | json_pp"
}

# Property CRUD
putProperty1="$(cat <<'END'
curl -fsS -X PUT 'localhost:8080/property/property1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{
  "street": "123 Main st.",
  "city": "Dallas",
  "state": "TX",
  "zip": "75401"
}' | json_pp
END
)"
getProperty1="$(cat <<'END'
curl -fsS -X GET 'localhost:8080/property/property1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
END
)"
putProperty2="$(cat <<'END'
curl -fsS -X PUT 'localhost:8080/property/property2' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{
  "street": "124 Main st.",
  "city": "Dallas",
  "state": "TX",
  "zip": "75401"
}' | json_pp
END
)"
delProperty2="$(cat <<'END'
curl -fsS -X DELETE 'localhost:8080/property/property2' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret'
END
)"
getProperties="$(cat <<'END'
curl -fsS -X GET 'localhost:8080/property' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
END
)"

main "$@"