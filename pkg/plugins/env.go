package plugins

import (
	"os"
)

const (
	AutobotClusterID   = "AUTOBOT_CLUSTER_ID"
	AutobotClusterURL  = "AUTOBOT_CLUSTER_URL"
	AutobotTopicEvents = "AUTOBOT_TOPIC_EVENTS"
)

// Env ...
type Env struct {
	env map[string]string
}

// DefaultEnv
func DefaultEnv() *Env {
	e := new(Env)
	e.env = map[string]string{
		AutobotClusterID:   "",
		AutobotClusterURL:  "",
		AutobotTopicEvents: "",
	}

	configureEnv(e)

	return e
}

// Get ...
func (e *Env) Get(s string) string {
	return e.env[s]
}

func configureEnv(e *Env) {
	for k, _ := range e.env {
		e.env[k] = os.Getenv(k)
	}
}
