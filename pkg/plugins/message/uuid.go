package message

import (
	"github.com/google/uuid"
)

// NewUUID ...
func NewUUID() string {
	return uuid.New().String()
}
