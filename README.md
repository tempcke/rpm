[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]

# Rental Property Manager (RPM)

[OpenAPI docs](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/tempcke/rpm/master/api/rest/openapi/openapi.yml)

run `make help` for a list of commands
running the tests requires some initialization which `make test` will do for you.

## Implemented Features both REST and gRPC
- **Property**: Store, Get, Lst, Remove
- **Tenant**:   Store, Get, List

## Roadmap
- Property Search API
- List Properties
    - filter
    - sort
    - paginate
- Lease property support
    - tenants
    - rent amount
    - term
    - payments
- property maintenance
    - ticket tracking
    - contractors
- Command Line Interface (CLI)
    
## RESTful API requests
### Add Property
```
curl -X POST "http://localhost:8080/property" \
     -H 'Content-Type: application/json' \
     -H 'Accept: application/json' \
     -d '{
  "street": "123 Main st.",
  "city": "Dallas",
  "state": "TX",
  "zip": "75401"
}' | json_pp
```

### List Properties
```
curl -X GET "http://localhost:8080/property" \
     -H 'Accept: application/json' | json_pp
```

### Get Property Detail
```
curl -X GET "http://localhost:8080/property/{propertyId}" \
     -H 'Accept: application/json' | json_pp
```

### Delete Property
```
curl -X DELETE "http://localhost:8080/property/{propertyId}"
```

[build-img]: https://github.com/tempcke/rpm/actions/workflows/test.yml/badge.svg
[build-url]: https://github.com/tempcke/rpm/actions
[pkg-img]: https://pkg.go.dev/badge/tempcke/rpm
[pkg-url]: https://pkg.go.dev/github.com/tempcke/rpm
[reportcard-img]: https://goreportcard.com/badge/tempcke/rpm
[reportcard-url]: https://goreportcard.com/report/tempcke/rpm