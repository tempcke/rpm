package config

// ENV keys
const (
	RequestID = "REQUEST_ID"
	LogLevel  = "LOG_LEVEL"
	AppName   = "APP_NAME"
	AppPort   = "APP_PORT"
	GrpcPort  = "GRPC_PORT"

	PostgresDSN           = "POSTGRES_DSN"
	PostgresDB            = "POSTGRES_DB"
	PostgresHost          = "POSTGRES_HOST"
	PostgresHostReplica   = "POSTGRES_READ_REPLICA_HOST"
	PostgresPort          = "POSTGRES_PORT"
	PostgresUser          = "POSTGRES_USER"
	PostgresPass          = "POSTGRES_PASSWORD"
	PostgresSSL           = "POSTGRES_SSLMODE"
	PostgresCertFile      = "POSTGRES_CERT_FILE" // location of cert
	PostgresMaxConnect    = "POSTGRES_MAX_CONNECTIONS"
	PostgresWritePoolOnly = "POSTGRES_WRITE_POOL_ONLY"
)
