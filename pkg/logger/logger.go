package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = log.With().
		Timestamp().
		Caller().
		Logger()
}

func Info(msg string, fields map[string]interface{}) {
	event := log.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func Debug(msg string, fields map[string]interface{}) {
	event := log.Debug()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func Warn(msg string, fields map[string]interface{}) {
	event := log.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}

func Error(msg string, err error, fields map[string]interface{}) {
	event := log.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}

	event.Msg(msg)
}

func Fatal(msg string, err error, fields map[string]interface{}) {
	event := log.Fatal().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}
