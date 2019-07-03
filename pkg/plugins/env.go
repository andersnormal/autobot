package plugins

import (
	"os"
	"strings"
)

var defaultEnv = map[string]string{
	AutobotName:          "",
	AutobotClusterID:     "",
	AutobotClusterURL:    "",
	AutobotTopicMessages: "",
	AutobotTopicReplies:  "",
}

// Env ...
type Env map[string]string

// NewDefaultEnv ...
func NewDefaultEnv() Env {
	env := copy(defaultEnv)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")

		env[pair[0]] = pair[1]
	}

	return env
}

// AutobotName ...
func (e Env) AutobotName() string {
	return e[AutobotName]
}

// AutobotClusterID ...
func (e Env) AutobotClusterID() string {
	return e[AutobotClusterID]
}

// AutobotClusterURL ...
func (e Env) AutobotClusterURL() string {
	return e[AutobotClusterURL]
}

// AutobotTopicMessages ...
func (e Env) AutobotTopicMessages() string {
	return e[AutobotTopicMessages]
}

// AutobotTopicReplies ...
func (e Env) AutobotTopicReplies() string {
	return e[AutobotTopicReplies]
}

// Getenv ...
func (e Env) Getenv(s string) string {
	return e[s]
}

func copy(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}
