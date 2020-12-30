# Rental Property Manager (RPM)

## Implemented Features
- Add Property
- List Properties
- Get Property Details
- Delete Property

## Roadmap
- Docker
    - PostgresSQL db and repository
    - .env for db credentials etc
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