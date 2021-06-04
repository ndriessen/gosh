package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

var log zerolog.Logger

func Trace(msg ...interface{}) {
	log.Trace().Msg(fmt.Sprint(msg...))
}

func Tracef(msg string, args ...interface{}) {
	log.Trace().Msgf(msg, args...)
}

func Debug(msg ...interface{}) {
	log.Debug().Msg(fmt.Sprint(msg...))
}

func Debugf(msg string, args ...interface{}) {
	log.Debug().Msgf(msg, args...)
}

func Info(msg ...interface{}) {
	log.Info().Msg(fmt.Sprint(msg...))
}

func Infof(msg string, args ...interface{}) {
	log.Info().Msgf(msg, args...)
}

func Warn(msg ...interface{}) {
	log.Warn().Msg(fmt.Sprint(msg...))
}

func Warnf(msg string, args ...interface{}) {
	log.Warn().Msgf(msg, args...)
}

func Err(err error, msg ...interface{}) error {
	log.Error().Err(err).Msg(fmt.Sprint(msg...))
	return err
}

func Errf(err error, msg string, args ...interface{}) error {
	log.Error().Err(err).Msgf(msg, args...)
	return err
}

func CheckErr(err error, args ...interface{}) error {
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprint(args...))
	}
	return err
}

func Fatal(err error, msg string, args ...interface{}) {
	log.Fatal().Err(err).Msgf(msg, args...)
	os.Exit(1)
}

func SetDetailedLogging() {
	output.PartsExclude = nil
	log = zerolog.New(output).With().Caller().Logger()
}

var output zerolog.ConsoleWriter

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func init() {
	zerolog.SetGlobalLevel(zerolog.FatalLevel)
	zerolog.CallerSkipFrameCount = 3
	output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen}
	output.PartsExclude = []string{"timestamp", "caller", "field", "level"}
	output.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	output.FormatLevel = func(i interface{}) string {
		level := strings.ToLower(i.(string))
		color := colorWhite
		switch level {
		case "info":
			color = colorGreen
		case "error", "fatal":
			color = colorRed
		case "warn":
			color = colorYellow
		}
		return fmt.Sprintf("\u001B[%dm%-6s\u001B[0m |", color, strings.ToUpper(level))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	output.FormatCaller = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%-60s |", i))
	}
	log = zerolog.New(output).With().Caller().Logger()
}
