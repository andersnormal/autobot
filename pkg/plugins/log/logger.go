package log

import (
	log "github.com/sirupsen/logrus"
)

// InternalLogger is the internal logging wrapper for the nlopes lib.
type InternalLogger struct {
	Logger *log.Entry
}

// Output ...
func (l *InternalLogger) Output(depth int, msg string) error {
	l.Logger.Info(msg)

	return nil
}
