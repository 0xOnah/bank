package logger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/0xOnah/bank/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type ContextKey string

const CtxKey ContextKey = "contextkey"

const (
	Production  = "production"
	Development = "development"
)

func SetUpLogger(level zerolog.Level, environment string) (*zerolog.Logger, error) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil, fmt.Errorf("failed to return buildinfo")
	}

	var output io.Writer
	switch environment {
	case Development:
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	case Production:
		output = os.Stdout
	default:
		return nil, errors.New("invalid enviroment variable")
	}

	log := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Int("pid", os.Getpid()).
		Str("environment", environment).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	return &log, nil
}

func InitLogger(cfg *config.Config) (*zerolog.Logger, error) {
	level := cfg.LEVEL
	if level == "" {
		level = "debug"
	}

	level = strings.TrimSpace(level)
	var logLevel zerolog.Level

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.ErrorLevel
	default:
		return nil, errors.New("invalid log level")
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger, err := SetUpLogger(logLevel, cfg.ENVIRONMENT)
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	return logger, nil
}

// func WithCTXLogger(ctx context.Context, logger *zerolog.Logger) *zerolog.Logger {
// 	return zerolog.Ctx(context.WithValue(ctx, CtxKey, logger))
// }

// func FromLoggerCTX(key string)
func ServiceLogger(log *zerolog.Logger, svcName string) *zerolog.Logger {
	svclog := log.With().Str("service", svcName).Logger()
	return &svclog
}
