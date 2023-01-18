
# opp-inventory API check
This is a generated document, to re-generate run: make apiCheck

## Property CRUD

### PUT property1
```
curl -fsS -X PUT 'localhost:8080/property/property1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{
  "street": "123 Main st.",
  "city": "Dallas",
  "state": "TX",
  "zip": "75401"
}' | json_pp | json_pp

ERROR!
```

### GET property1
```
curl -fsS -X GET 'localhost:8080/property/property1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp | json_pp

ERROR!
```

### PUT and DELETE property2
```
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

ERROR!
curl -fsS -X DELETE 'localhost:8080/property/property2' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret'

ERROR!
```

### GET properties
```
curl -fsS -X GET 'localhost:8080/property' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp | json_pp

ERROR!
```
