package log_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tempcke/rpm/internal/lib/log"
)

type myLogger struct {
	*logrus.Logger
	someInt int
}

func newMyLogger() myLogger {
	return myLogger{someInt: 42}
}

func TestLogger(t *testing.T) {
	t.Run("get default logger without setup", func(t *testing.T) {
		level := logrus.DebugLevel
		t.Setenv("LOG_LEVEL", level.String())
		l := log.Logger()
		require.NotNil(t, l)
		assert.Equal(t, level.String(), log.Entry().Logger.Level.String())
	})

	// be aware that
	t.Run("UseLogger", func(t *testing.T) {
		log.UseLogger(newMyLogger())
		logger := log.Logger()
		ml, ok := logger.(myLogger)
		require.True(t, ok, "UseLogger didn't give us a myLogger")
		assert.Equal(t, 42, ml.someInt)
	})
}
