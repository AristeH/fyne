package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var log = logrus.New()

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		TimestampFormat: "02/01/06 15:04:05",
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	err := os.MkdirAll("logs", 0644)
	if err != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(logrus.WarnLevel)
	logrus.SetOutput(io.Discard)
	log.AddHook(&writerHook{
		Writer:    []io.Writer{allFile},
		LogLevels: logrus.AllLevels,
	})
}

func GetLog() *logrus.Logger {
	return log
}
