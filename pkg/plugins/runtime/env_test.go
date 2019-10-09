package runtime

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnInitialize(t *testing.T) {
	assert := assert.New(t)

	run := false

	OnInitialize(func() { run = true })

	runInitializers()

	assert.True(run)
	assert.Len(initializers, 1)
}

func TestNewRuntime(t *testing.T) {
	assert := assert.New(t)

	runtime := NewRuntime()

	assert.NotNil(runtime)
}

func TestRuntime_hasRuntimeFunc(t *testing.T) {
	tests := []struct {
		desc     string
		runtime  *Runtime
		expected bool
	}{
		{
			desc:     "missing run funcs",
			runtime:  &Runtime{},
			expected: false,
		},
		{
			desc: "has run func but no run with error func",
			runtime: &Runtime{
				Run: func(*Environment) {},
			},
			expected: true,
		},
		{
			desc: "has run error func but no run func",
			runtime: &Runtime{
				RunE: func(*Environment) error { return nil },
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			assert.NotNil(tt.runtime)
			assert.Equal(tt.expected, tt.runtime.hasRuntimeFuncs())
		})
	}
}

func TestRuntimeExecute(t *testing.T) {
	tests := []struct {
		desc           string
		env            *Environment
		expectedEnv    *Environment
		runtime        *Runtime
		expectedRunErr error
		err            error
	}{
		{
			desc:           "missing run funcs",
			runtime:        &Runtime{},
			env:            nil,
			expectedEnv:    nil,
			expectedRunErr: nil,
			err:            ErrNoRuntimeFunc,
		},
		{
			desc: "has run func",
			runtime: &Runtime{Run: func(env *Environment) {
				env.Debug = true
			}},
			expectedRunErr: nil,
			env:            &Environment{},
			expectedEnv: &Environment{
				Debug: true,
			},
			err: nil,
		},
		{
			desc: "has run func with no error",
			runtime: &Runtime{RunE: func(env *Environment) error {
				env.Debug = true

				return nil
			}},
			expectedRunErr: nil,
			env:            &Environment{},
			expectedEnv: &Environment{
				Debug: true,
			},
			err: nil,
		},
		{
			desc: "has run func with error",
			runtime: &Runtime{RunE: func(env *Environment) error {
				return errors.New("foo")
			}},
			env:            &Environment{},
			expectedEnv:    &Environment{},
			err:            nil,
			expectedRunErr: errors.New("foo"),
		},
		// {
		// 	desc: "has run func",
		// 	runtime: &Runtime{
		// 		Run: func(*Environment) {},
		// 	},
		// 	expected: true,
		// },
		// {
		// 	desc: "has run error func",
		// 	runtime: &Runtime{
		// 		RunE: func(*Environment) error { return nil },
		// 	},
		// 	expected: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert := assert.New(t)

			if tt.err != nil {
				err := tt.runtime.Execute()

				assert.Error(err)
				assert.EqualError(tt.err, err.Error())

				return
			}

			env = tt.env
			err := tt.runtime.Execute()

			if tt.expectedRunErr != nil {
				assert.Error(err)
				assert.EqualError(tt.expectedRunErr, err.Error())

				return
			}

			assert.NoError(err)
			assert.Equal(tt.expectedEnv, tt.env)
		})
	}
}
