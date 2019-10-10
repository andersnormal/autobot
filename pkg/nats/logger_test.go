package nats

import (
	"bytes"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	logger := NewLogger()

	assert.NotNil(logger)
}

func TestLogger_SetLogger(t *testing.T) {
	assert := assert.New(t)

	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)
}

func TestLogger_Errorf(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)

	logger.Errorf("%s %s", "foo", "bar")
	assert.Contains(buf.String(), "foo bar")
	assert.Contains(buf.String(), "level=error")

	log.SetOutput(os.Stderr)
}

func TestLogger_Debugf(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.DebugLevel)

	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)

	logger.Debugf("%s %s", "foo", "bar")
	assert.Contains(buf.String(), "foo bar")
	assert.Contains(buf.String(), "level=debug")

	log.SetOutput(os.Stderr)
}

func TestLogger_Noticef(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.InfoLevel)

	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)

	logger.Noticef("%s %s", "foo", "bar")
	assert.Contains(buf.String(), "foo bar")
	assert.Contains(buf.String(), "level=info")

	log.SetOutput(os.Stderr)
}

func TestLogger_Warnf(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.WarnLevel)

	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)

	logger.Warnf("%s %s", "foo", "bar")
	assert.Contains(buf.String(), "foo bar")
	assert.Contains(buf.String(), "level=warn")

	log.SetOutput(os.Stderr)
}

func TestLogger_Tracef(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.TraceLevel)

	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)

	logger.Tracef("%s %s", "foo", "bar")
	assert.Contains(buf.String(), "foo bar")
	assert.Contains(buf.String(), "level=trace")

	log.SetOutput(os.Stderr)
}

func TestLogger_Fatalf(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetLevel(log.ErrorLevel)

	logger := NewLogger()
	assert.NotNil(logger)

	logger.SetLogger(log.WithFields(log.Fields{}))
	assert.NotNil(logger.log)

	log.SetOutput(os.Stderr)
}
