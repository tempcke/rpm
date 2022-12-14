package config_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal/config"
)

func TestConfigBuilder(t *testing.T) {
	t.Run("ImportMap", func(t *testing.T) {
		cm1 := config.ConfMap{
			config.LogLevel:  logrus.InfoLevel.String(),
			config.RequestID: "foo",
		}

		// just to demonstrate you can import multiple maps
		requestID := "foo-request"
		cm2 := config.ConfMap{config.RequestID: requestID}

		c, err := config.NewConfBuilder().
			ImportMap(cm1).
			ImportMap(cm2).
			Build()

		require.NoError(t, err)
		assert.Equal(t, logrus.InfoLevel.String(), c.LogLevel())
		assert.Equal(t, requestID, c.RequestID())
	})

	t.Run("ImportFile", func(t *testing.T) {
		file, cleanupFunc := tempEnvFile(t)
		defer cleanupFunc()

		cm := config.ConfMap{
			config.LogLevel: logrus.WarnLevel.String(),
		}

		writeEnvFile(t, file, cm)

		c, err := config.NewConfBuilder().
			ImportFile(file).
			Build()

		require.NoError(t, err)
		assert.Equal(t, logrus.WarnLevel.String(), c.LogLevel())
	})

	t.Run("AutomaticEnv", func(t *testing.T) {
		pgHost := "some-host"

		// set os env but preserve system state
		origVal := os.Getenv(config.LogLevel)
		defer func() { _ = os.Setenv(config.PostgresHost, origVal) }()
		_ = os.Setenv(config.PostgresHost, pgHost)

		c, err := config.NewConfBuilder().
			AutomaticEnv().
			ImportMap(config.ConfMap{config.PostgresHost: "bar"}).
			Build()

		require.NoError(t, err)
		assert.Equal(t, pgHost, c.PostgresHost())
	})

	t.Run("AutomaticEnv over confMap", func(t *testing.T) {
		cm := config.ConfMap{
			config.LogLevel: logrus.WarnLevel.String(),
		}

		// set os env but preserve system state
		originalLogLevel := os.Getenv(config.LogLevel)
		defer func() { _ = os.Setenv(config.LogLevel, originalLogLevel) }()
		_ = os.Setenv(config.LogLevel, logrus.ErrorLevel.String())

		c, err := config.NewConfBuilder().
			AutomaticEnv().
			ImportMap(cm).
			Build()

		require.NoError(t, err)
		assert.Equal(t, logrus.ErrorLevel.String(), c.LogLevel())
	})

	t.Run("can use viper.GetString() to fetch values", func(t *testing.T) {
		requestID := uuid.NewString()
		cm := config.ConfMap{config.RequestID: requestID}

		c, err := config.NewConfBuilder().
			ImportMap(cm).
			Build()

		require.NoError(t, err)
		assert.Equal(t, requestID, c.RequestID())
		assert.NotEqual(t, requestID, viper.GetString(config.RequestID))
	})
	t.Run("can viper get these on it's own?", func(t *testing.T) {
		requestID := uuid.NewString()
		cm := config.ConfMap{config.RequestID: requestID}

		c, err := config.NewConfBuilderGlobal().
			ImportMap(cm).
			Build()

		require.NoError(t, err)
		assert.Equal(t, requestID, c.RequestID())
		assert.Equal(t, requestID, viper.GetString(config.RequestID))
	})
}

func TestConfig(t *testing.T) {
	var (
		logLevel = logrus.InfoLevel.String()
		appName  = "some-appName"
	)
	c, err := config.NewConfBuilder().
		ImportMap(config.ConfMap{
			config.LogLevel: logrus.InfoLevel.String(),
			config.AppName:  appName,
		}).
		Build()
	require.NoError(t, err)

	t.Run("helper methods to get types", func(t *testing.T) {
		assert.Equal(t, logLevel, c.GetLogLevel().String())
		assert.Equal(t, appName, c.AppName())
	})
}

func tempEnvFile(t *testing.T) (file string, cleanupFunc func()) {
	dir, err := os.MkdirTemp("", "google-tdd-test-")
	require.NoError(t, err)
	file = fmt.Sprintf("%s/app.env", dir)
	cleanupFunc = func() { _ = os.RemoveAll(dir) }
	return file, cleanupFunc
}

func writeEnvFile(t *testing.T, file string, cm config.ConfMap) {
	t.Helper()
	var sb strings.Builder
	for k, v := range cm {
		sb.WriteString(fmt.Sprintf("%v=%v\n", k, v))
	}
	err := os.WriteFile(file, []byte(sb.String()), 0666)
	require.NoError(t, err)
}
