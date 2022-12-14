package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	defaultLogLevel = logrus.InfoLevel.String()
)

// ConfMap defines a key value map to be imported into the config
// see config/key package for valid keys to use
type ConfMap map[string]interface{}

// Config is the application configuration
// with some getters to help get things in more desirable types
type Config struct {
	v      *viper.Viper
	logger logrus.FieldLogger
}

func (c Config) WithLogger(l logrus.FieldLogger) Config {
	c.logger = l
	return c
}

func (c Config) viper() *viper.Viper {
	if c.v == nil {
		return viper.GetViper()
	}
	return c.v
}

func (c Config) GetString(key string) string     { return c.viper().GetString(key) }
func (c Config) GetInt(key string) int           { return c.viper().GetInt(key) }
func (c Config) GetInt64(key string) int64       { return c.viper().GetInt64(key) }
func (c Config) GetFloat64(key string) float64   { return c.viper().GetFloat64(key) }
func (c Config) GetBool(key string) bool         { return c.viper().GetBool(key) } // see strconv.ParseBool
func (c Config) Get(key string) interface{}      { return c.viper().Get(key) }
func (c Config) Set(key string, val interface{}) { c.viper().Set(key, val) }

func (c Config) AppName() string   { return c.GetString(AppName) }
func (c Config) LogLevel() string  { return c.GetString(LogLevel) }
func (c Config) RequestID() string { return c.GetString(RequestID) }

func (c Config) PostgresDSN() string         { return c.GetString(PostgresDSN) }
func (c Config) PostgresDB() string          { return c.GetString(PostgresDB) }
func (c Config) PostgresHost() string        { return c.GetString(PostgresHost) }
func (c Config) PostgresHostReplica() string { return c.GetString(PostgresHostReplica) }
func (c Config) PostgresPort() int           { return c.GetInt(PostgresPort) }
func (c Config) PostgresSSL() string         { return c.GetString(PostgresSSL) }
func (c Config) PostgresUser() string        { return c.GetString(PostgresUser) }
func (c Config) PostgresPass() string        { return c.GetString(PostgresPass) }
func (c Config) PostgresMaxConnect() int     { return c.GetInt(PostgresMaxConnect) }
func (c Config) PostgresWritePoolOnly() bool { return c.GetBool(PostgresWritePoolOnly) }
func (c Config) PostgresCertFile() string {
	// when path is set, just return it
	if filePath := c.GetString(PostgresCertFile); filePath != "" {
		return filePath
	}

	// else there simply is none
	return ""
}
func (c Config) GetPostgresMaxConnect(fallback int) int {
	if c.PostgresMaxConnect() < 1 {
		return fallback
	}
	return c.PostgresMaxConnect()
}

// GetLogLevel converts the string loglevel into a real logrus.Level
func (c Config) GetLogLevel() logrus.Level {
	lvl, err := logrus.ParseLevel(c.LogLevel())
	if err != nil {
		lvl, _ = logrus.ParseLevel(defaultLogLevel)
	}
	return lvl
}
