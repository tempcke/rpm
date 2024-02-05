
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
{"city":"Dallas","id":"property2","state":"TX","street":"124 Main st.","zip":"75401"}
curl -fsS -X DELETE 'localhost:8080/property/property2' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret'

```

### GET properties
```
curl -fsS -X GET 'localhost:8080/property' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
{
   "properties" : [
      {
         "city" : "Dallas",
         "id" : "property1",
         "state" : "TX",
         "street" : "123 Main st.",
         "zip" : "75401"
      }
   ]
}
```

### Search properties
```
curl -fsS -X GET 'localhost:8080/property?search=dallas%20tx' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
{
   "properties" : [
      {
         "city" : "Dallas",
         "id" : "property1",
         "state" : "TX",
         "street" : "123 Main st.",
         "zip" : "75401"
      }
   ]
}
```

## Tenant CRUD

### PUT tenant1
```
curl -fsS -X PUT 'localhost:8080/tenant/tenant1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{
  "tenant": {
    "fullName": "John Doe",
    "dlNum": "646673153",
    "dlState": "TX",
    "dob": "2006-01-02",
    "phones": [{"number": "555-555-1234", "desc": "mobile"}]
  }
}' | json_pp
{
   "tenant" : {
      "dlNum" : "646673153",
      "dlState" : "TX",
      "dob" : "2006-01-02",
      "fullName" : "John Doe",
      "id" : "tenant1",
      "phones" : [
         {
            "desc" : "mobile",
            "number" : "555-555-1234"
         }
      ]
   }
}
```

### GET tenant1
```
curl -fsS -X GET 'localhost:8080/tenant/tenant1' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
{
   "tenant" : {
      "dlNum" : "646673153",
      "dlState" : "TX",
      "dob" : "2006-01-02",
      "fullName" : "John Doe",
      "id" : "tenant1",
      "phones" : [
         {
            "desc" : "mobile",
            "number" : "555-555-1234"
         }
      ]
   }
}
```

### POST tenant
```
curl -fsS -X POST 'localhost:8080/tenant' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -d '{
  "tenant": {
    "fullName": "Jane Doe",
    "dlNum": "746673153",
    "dlState": "TX",
    "dob": "2006-01-02",
    "phones": [{"number": "555-555-1235", "desc": "mobile"}]
  }
}' | json_pp
{
   "tenant" : {
      "dlNum" : "746673153",
      "dlState" : "TX",
      "dob" : "2006-01-02",
      "fullName" : "Jane Doe",
      "id" : "b8eb9b54-ccac-4dfd-bc7b-6e6da5925a60",
      "phones" : [
         {
            "desc" : "mobile",
            "number" : "555-555-1235"
         }
      ]
   }
}
```

### GET tenants
```
curl -fsS -X GET 'localhost:8080/tenant' \
  -H 'X-API-Key: key' -H 'X-API-Secret: secret' | json_pp
{
   "tenants" : [
      {
         "dlNum" : "646673153",
         "dlState" : "TX",
         "dob" : "2006-01-02",
         "fullName" : "John Doe",
         "id" : "tenant1",
         "phones" : null
      },
      {
         "dlNum" : "746673153",
         "dlState" : "TX",
         "dob" : "2006-01-02",
         "fullName" : "Jane Doe",
         "id" : "b8eb9b54-ccac-4dfd-bc7b-6e6da5925a60",
         "phones" : null
      }
   ]
}
```
