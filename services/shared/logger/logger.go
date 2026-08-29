package logger

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func New(w io.Writer) Logger {
	consoleWriter := zerolog.ConsoleWriter{Out: w, TimeFormat: zerolog.TimeFormatUnix, TimeLocation: time.Local}
	multi := zerolog.MultiLevelWriter(consoleWriter, os.Stdout)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	return Logger{logger: logger}
}

func (l Logger) PrintInfo(message string, properties map[string]string) {
	json, err := json.Marshal(properties)
	if err != nil {
		log.Printf("Info: %v %s %+v", time.Now(), message, properties)
		return
	}
	l.logger.Info().RawJSON("properties", json).Msg(message)
}

func (l Logger) PrintError(err error, message string, properties map[string]string) {
	json, err := json.Marshal(properties)
	if err != nil {
		log.Printf("Info: %v %s %s %+v", time.Now(), err.Error(), message, properties)
		return
	}
	logger := l.logger.With().Stack().Logger()

	logger.Error().Err(err).RawJSON("properties", json).Msg(message)
}

func (l Logger) PrintFatal(err error, message string, properties map[string]string) {
	json, err := json.Marshal(properties)
	if err != nil {
		log.Printf("Info: %v %s %s %+v", time.Now(), err.Error(), message, properties)
		return
	}
	logger := l.logger.With().Stack().Logger()

	logger.Fatal().Err(err).RawJSON("properties", json).Msg(message)
}

func (l Logger) PrintWarn(err error, message string, properties map[string]string) {
	json, err := json.Marshal(properties)
	if err != nil {
		log.Printf("Info: %v %s %s %+v", time.Now(), err.Error(), message, properties)
		return
	}
	logger := l.logger.With().Stack().Logger()

	logger.Warn().Err(err).RawJSON("properties", json).Msg(message)
}

func (l Logger) PrintPanic(err error, message string, properties map[string]string) {
	json, err := json.Marshal(properties)
	if err != nil {
		log.Printf("Info: %v %s %s %+v", time.Now(), err.Error(), message, properties)
		return
	}
	logger := l.logger.With().Stack().Logger()

	logger.Panic().Err(err).RawJSON("properties", json).Msg(message)
}
