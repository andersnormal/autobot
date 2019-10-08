package message_test

import (
	"testing"

	. "github.com/andersnormal/autobot/pkg/plugins/message"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	assert := assert.New(t)

	uuid := NewUUID()

	assert.NotEmpty(uuid)
	assert.Len(uuid, len("00000000-0000-0000-0000-000000000000"))
}
