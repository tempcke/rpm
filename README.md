[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]

# Rental Property Manager (RPM)

**WIP**: Be advised that this code is NOT used in production anywhere and this project is really more of a portfolio project than something that is likely to ever be used.

- [OpenAPI docs](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/tempcke/rpm/master/api/rest/openapi/openapi.yml)
- [generated api docs](apicheck.md)

## make help
run `make help` for a list of commands
running the tests requires some initialization which `make test` will do for you.

```
Usage: make <target>
  help              display help information
  build             create binary
  run               create and run binary
  check             lint + test, pre-commit hook
  lint              fmt, vet, and staticcheck
  test              execute tests
  testAll           run all tests including those that need docker/postgres
  testCI            exact tests the way buildkite does, use for local debug of buildkite failure
  testAcceptance    black box testing
  dockerUp          docker-compose up
  dockerDown        docker-compose down
  dockerRestart     dockerDown && dockerUp
  dockerRestartApp  rebuild app and replace running container
  dockerFollowLogs  live stream logs from docker-compose
  apiCheck          generate api docs in apicheck.md
  clean             dockerDown && docker-compose down for CI
  oapigen           generate api/openapi/*.gen.go using github.com/deepmap/oapi-codegen
  protoc            generate api/rpc/proto/*.pb.go
  cert              Create certificates to encrypt the gRPC connection
```

## Design concepts
- **TDD** (Red-Green-Refactor) as described by Kent Beck and Robert "Uncle Bob" Martin
  - Tests focus on the behavior of the system and must avoid the implementation details.
  - Not every detail is tested at every layer.
  - Core logic must be driven by failing tests.
  - Prefer use-case layer tests against a real postgres database running in docker
    - This avoids having to decide what to test at the repository later and what to test in the use-case layer.  The use-case describes the desired behavior of the system.  So test it there.
  - Use repository layer testing to drive error handling scenarios and other edge cases
  - Service/API layer testing should be minimal just to ensure routing is set up correctly and request/responses are as desired.
- **Clean Architecture**
  - AKA: Hexagonal Architecture, Ports & Adapters, Onion Architecture, etc.
  - Maintain separation of concerns at different layers.  Ports -> Adapters -> Use Cases -> Entities
  - Thanks [UncleBob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- **Multiple APIs backed by the same business logic**
  - Each action is exposed in a similar way by both **GRPC** and **REST** APIs
  - All api actions are implemented in a use-case layer.  So there is one place to implement and test the core use cases while both grpc and rest can expose it.
- **OpenAPI**
  - Having up-to-date API docs is important but can be challenging.
  - To change the RESTful API [openapi.yaml](api/rest/openapi/openapi.yml) is updated and then `make oapigen` should be run which generates the go-chi router, server interface, and api models.
  - Thanks [Three Dots Labs](https://threedots.tech/post/list-of-recommended-libraries/#openapi)
- **Model Separation**
  - API and domain models must be separate.  They often start off the same or very similar, but they change for different reasons.
  - API models are not used outside /api
    - There are functions to map to and from domain models under the api package.
  - Domain models are used everywhere else and found in /entity.
    - This allows the use-case layer to use just one set of models and the grpc and rest services map to/from them in lightweight handlers.
- **Lightweight Handlers**
  - Handlers must be lightweight.  Read the request, send to the use-case layer, and then compose the response, that is it.
- **Acceptance testing**
  - Use specifications which take a driver to test the basic flows through the system against a compiled binary.
    - This tests the entire system built by main.
    - Should be very few tests, only enough to know that main and the servers are stood up and wired correctly.
    - Obviously the server should be backed by a real postgres database, not in-memory
    - tests do not have direct access to the database and MUST go through the API for everything
  - test specifications use a driver interface.  Other layers can implement the driver and be tested through the spec also.  This allows for much easier debugging when a specification test fails for a given service.
  - Thanks [quii](https://quii.gitbook.io/learn-go-with-tests/testing-fundamentals/scaling-acceptance-tests#separation-of-concerns) and [Dave Farley](https://www.youtube.com/watch?v=JDD5EEJgpHU)
## Implemented Features both REST and gRPC
- **Property**: 
  - Store, Get, Remove
  - List with search string filter
- **Tenant**:
  - Store, Get, List

## Roadmap
- filter, sort, paginate
- Prometheus
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
There is a shell script [apicheck.sh](apicheck.sh) which can be executed via `make apiCheck`.  This script calls the restful endpoints on the binary application running in docker and generates [apicheck.md](apicheck.md).  The goal of this is to auto generate api example docs from hitting the actual running api, so you can see what the real requests and responses look like.

[build-img]: https://github.com/tempcke/rpm/actions/workflows/test.yml/badge.svg
[build-url]: https://github.com/tempcke/rpm/actions
[pkg-img]: https://pkg.go.dev/badge/tempcke/rpm
[pkg-url]: https://pkg.go.dev/github.com/tempcke/rpm
[reportcard-img]: https://goreportcard.com/badge/tempcke/rpm
[reportcard-url]: https://goreportcard.com/report/tempcke/rpm