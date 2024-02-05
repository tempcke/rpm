[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]

# Rental Property Manager (RPM)

**WIP**: Be advised that this code is NOT used in production anywhere and this project is really more of a portfolio project than something that is likely to ever be used.

- [OpenAPI docs](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/tempcke/rpm/master/api/rest/openapi/openapi.yml)
- [generated api docs](https://github.com/tempcke/rpm/blob/master/apicheck.md)

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
There is a shell script [apicheck.sh](https://github.com/tempcke/rpm/blob/master/apicheck.sh) which can be executed via `make apiCheck`.  This script calls the restful endpoints on the binary application running in docker and generates [apicheck.md](https://github.com/tempcke/rpm/blob/master/apicheck.md).  The goal of this is to auto generate api example docs from hitting the actual running api, so you can see what the real requests and responses look like.

[build-img]: https://github.com/tempcke/rpm/actions/workflows/test.yml/badge.svg
[build-url]: https://github.com/tempcke/rpm/actions
[pkg-img]: https://pkg.go.dev/badge/tempcke/rpm
[pkg-url]: https://pkg.go.dev/github.com/tempcke/rpm
[reportcard-img]: https://goreportcard.com/badge/tempcke/rpm
[reportcard-url]: https://goreportcard.com/report/tempcke/rpm