package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	assert := assert.New(t)

	err := NewError("%s %s", "foo", "bar")

	assert.NotNil(err)
	assert.NotEmpty(err)
	assert.Equal(err, queueError{err: "foo bar"})
}
