package log

import (
	"bytes"
	"os"
	"testing"

	ll "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInternalLogger(t *testing.T) {
	assert := assert.New(t)

	var buf bytes.Buffer
	ll.SetOutput(&buf)

	logger := &InternalLogger{
		Logger: ll.WithFields(ll.Fields{}),
	}

	err := logger.Output(1, "foo")

	ll.SetOutput(os.Stderr)

	assert.NoError(err)
	assert.NotEmpty(buf.String())
	assert.Contains(buf.String(), "foo")
}
