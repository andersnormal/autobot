package message_test

import (
	"fmt"
	"testing"

	. "github.com/andersnormal/autobot/pkg/plugins/message"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	tests := []struct {
		desc     string
		metadata Metadata
		key      string
		value    string
	}{
		{
			desc: "existing key should return string",
			metadata: Metadata{
				"foo": "bar",
			},
			key:   "foo",
			value: "bar",
		},
		{
			desc: "not existing key should return empty string",
			metadata: Metadata{
				"foo": "bar",
			},
			key:   "bar",
			value: "",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			a := assert.New(t)

			var v string
			a.NotPanics(func() { v = tt.metadata.Get(tt.key) })
			a.Equal(tt.value, v)
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		desc     string
		metadata Metadata
		key      string
		value    string
	}{
		{
			desc:     "should set a key to a value",
			metadata: Metadata{},
			key:      "foo",
			value:    "bar",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			a := assert.New(t)

			a.NotPanics(func() { tt.metadata.Set(tt.key, tt.value) })
			a.Equal(tt.metadata[tt.key], tt.value)
		})
	}
}

func TestSrc(t *testing.T) {
	tests := []struct {
		desc     string
		metadata Metadata
		value    string
	}{
		{
			desc:     "should set source",
			metadata: Metadata{},
			value:    "foo",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			a := assert.New(t)

			a.NotPanics(func() { tt.metadata.Src(tt.value) })
			a.Equal(tt.metadata[Src], tt.value)
		})
	}
}

func TestDest(t *testing.T) {
	tests := []struct {
		desc     string
		metadata Metadata
		value    string
	}{
		{
			desc:     "should set destination",
			metadata: Metadata{},
			value:    "bar",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.desc), func(t *testing.T) {
			a := assert.New(t)

			a.NotPanics(func() { tt.metadata.Dest(tt.value) })
			a.Equal(tt.metadata[Dest], tt.value)
		})
	}
}
