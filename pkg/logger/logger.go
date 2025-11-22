package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Logger zerolog.Logger

// Config holds logger configuration
type Config struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// Init initializes the global logger
func Init(cfg Config) error {
	// Set log level
	level, err := parseLogLevel(cfg.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)

	// Configure time format
	zerolog.TimeFieldFormat = time.RFC3339

	var writers []io.Writer

	// Console output
	if cfg.Output == "stdout" || cfg.Output == "both" {
		var consoleWriter io.Writer
		if cfg.Format == "pretty" {
			consoleWriter = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
				NoColor:    false,
			}
		} else {
			consoleWriter = os.Stdout
		}
		writers = append(writers, consoleWriter)
	}

	// File output
	if cfg.Output == "file" || cfg.Output == "both" {
		if cfg.FilePath != "" {
			// Ensure directory exists
			dir := filepath.Dir(cfg.FilePath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}

			file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			writers = append(writers, file)
		}
	}

	// Create multi-writer if multiple outputs
	var writer io.Writer
	if len(writers) > 1 {
		writer = io.MultiWriter(writers...)
	} else if len(writers) == 1 {
		writer = writers[0]
	} else {
		writer = os.Stdout
	}

	// Create logger
	Logger = zerolog.New(writer).
		With().
		Timestamp().
		Caller().
		Logger()

	// Set global logger
	log.Logger = Logger

	return nil
}

// parseLogLevel parses log level string
func parseLogLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn", "warning":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	case "panic":
		return zerolog.PanicLevel, nil
	default:
		return zerolog.InfoLevel, nil
	}
}

// WithRequestID returns a logger with request ID
func WithRequestID(requestID string) *zerolog.Logger {
	l := Logger.With().Str("request_id", requestID).Logger()
	return &l
}

// WithUserID returns a logger with user ID
func WithUserID(userID string) *zerolog.Logger {
	l := Logger.With().Str("user_id", userID).Logger()
	return &l
}

// WithContext returns a logger with custom context fields
func WithContext(fields map[string]interface{}) *zerolog.Logger {
	ctx := Logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	l := ctx.Logger()
	return &l
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Debug().Msg(msg)
}

// Info logs an info message
func Info(msg string) {
	Logger.Info().Msg(msg)
}

// Warn logs a warning message
func Warn(msg string) {
	Logger.Warn().Msg(msg)
}

// Error logs an error message
func Error(msg string, err error) {
	Logger.Error().Err(err).Msg(msg)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, err error) {
	Logger.Fatal().Err(err).Msg(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	Logger.Debug().Msgf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	Logger.Info().Msgf(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	Logger.Warn().Msgf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	Logger.Error().Msgf(format, args...)
}

// Get returns the global logger
func Get() *zerolog.Logger {
	return &Logger
}
