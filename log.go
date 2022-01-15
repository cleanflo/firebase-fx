package register

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.FatalLevel)
	log.SetFormatter(&log.JSONFormatter{
		DisableTimestamp: true,
	})
}

// SetLogLevel is a wrapper for setting the logrus.Level
func SetLogLevel(level LogLevel) {
	switch level {
	case Debug:
		log.SetLevel(log.DebugLevel)
	case Info:
		log.SetLevel(log.InfoLevel)
	case Warn:
		log.SetLevel(log.WarnLevel)
	case Error:
		log.SetLevel(log.ErrorLevel)
	}
}

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

// LogLevel is a helper to log at different levels
type LogLevel int

// Msgf is a wrapper for logrus.{Level}f
func (l LogLevel) Msgf(msg string, args ...interface{}) {
	switch l {
	case Debug:
		log.Debugf(msg, args...)
	case Info:
		log.Infof(msg, args...)
	case Warn:
		log.Warnf(msg, args...)
	case Error:
		log.Errorf(msg, args...)
	default:
		log.Debugf(msg, args...)
	}
}

// Msg is a wrapper for logrus.{Level}
func (l LogLevel) Msg(msg string, fields map[string]interface{}) {
	switch l {
	case Debug:
		log.Debug(msg, fields)
	case Info:
		log.Info(msg, fields)
	case Warn:
		log.Warn(msg, fields)
	case Error:
		log.Error(msg, fields)
	default:
		log.Debug(msg, fields)
	}
}

// Err is a wrapper for logrus.Error
func (l LogLevel) Err(msg string, err error) error {
	err = fmt.Errorf("%s: %v", msg, err)
	l.Msgf(msg, err)
	return err
}

// Errf is a wrapper for logrus.Errorf
func (l LogLevel) Errf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)
	l.Msgf("%s", err)
	return err
}
