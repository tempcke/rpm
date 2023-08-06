
# API check
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
}' | json_pp
{
   "city" : "Dallas",
   "createdAt" : "2023-08-05T20:21:25-05:00",
   "id" : "property1",
   "state" : "TX",
   "street" : "123 Main st.",
   "zip" : "75401"
}
```

### GET property1
```
curl -fsS -X GET 'localhost:8080/property/property1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
{
   "city" : "Dallas",
   "createdAt" : "2023-08-05T20:21:25-05:00",
   "id" : "property1",
   "state" : "TX",
   "street" : "123 Main st.",
   "zip" : "75401"
}
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
}'
{"id":"property2","street":"124 Main st.","city":"Dallas","state":"TX","zip":"75401","createdAt":"2023-08-05T20:21:25-05:00"}
curl -fsS -X DELETE 'localhost:8080/property/property2' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret'

```

### GET properties
```
curl -fsS -X GET 'localhost:8080/property' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
{
   "items" : [
      {
         "city" : "Dallas",
         "createdAt" : "2023-08-05T20:21:25-05:00",
         "id" : "property1",
         "state" : "TX",
         "street" : "123 Main st.",
         "zip" : "75401"
      }
   ]
}
```
