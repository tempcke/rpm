FROM golang:1.21-alpine AS builder

WORKDIR /app

# no longer needed because of vendor folder.
# go.mod go.sum
#COPY go.* ./
#RUN go mod download
#RUN go mod verify

# copy source files and build the binary
COPY . .
RUN go build -o bin/rpm cmd/rpmserver/*.go

FROM alpine:3.16.2
RUN apk add --no-cache \
    ca-certificates \
    postgresql-client

WORKDIR /app
COPY --from=builder /app/bin/rpm .
COPY ./service.pem .
COPY ./service.key .

CMD ["./rpm"]
