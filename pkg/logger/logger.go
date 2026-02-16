package logger

import (
	"fmt"
	"github.com/sirupsen/logrus/hooks/writer"
	"io"
	"os"

	"github.com/bdzhalalov/kolikosoft-trade/pkg/config"
	"github.com/sirupsen/logrus"
)

func Logger(config *config.Config) *logrus.Logger {

	log := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		fmt.Printf("Can't parse log level, setting level to DEBUG: %v\n", err.Error())
		level = logrus.DebugLevel
	}

	log.SetLevel(level)
	log.SetOutput(io.Discard)

	file, err := createLogFile()
	if err == nil {
		log.AddHook(&writer.Hook{
			Writer: file,
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
				logrus.DebugLevel,
			},
		})

		log.AddHook(&writer.Hook{
			Writer: os.Stdout,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
			},
		})
	} else {
		log.SetOutput(os.Stdout)
		log.Info("Failed to log to file, using default stdout")
	}

	return log
}

func createLogFile() (*os.File, error) {
	if err := os.MkdirAll("./logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
