package log

import (
	"log/slog"
)

type (
	SLogger interface {
		Debug(msg string, args ...any)
		Info(msg string, args ...any)
		Warn(msg string, args ...any)
		Error(msg string, args ...any)
	}
	Fields = map[string]any
)

func WithError(err error) *slog.Logger {
	return WithFields(Fields{"error": err})
}
func WithFields(fields ...Fields) *slog.Logger {
	entry := slog.Default()
	for _, f := range fields {
		for k, v := range f {
			entry = entry.With(k, v)
		}
	}
	return entry
}
func Entry(fields ...Fields) *slog.Logger {
	return WithFields(fields...)
}
func Error(msg string, args ...any) {
	Entry().Error(msg, args...)
}
