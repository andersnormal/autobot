package runtime

const (
	// DefaultClusterID ...
	DefaultClusterID = "autobot"
	// DefaultClusterURL ...
	DefaultClusterURL = "nats://localhost:4222"
	// DefaultClusterInbox ...
	DefaultClusterInbox = "autobot.inbox"
	// DefaultClusterOutbox ...
	DefaultClusterOutbox = "autobot.outbox"
	// DefaultLogFormat ...
	DefaultLogFormat = "text"
	// DefaultLogLevel ...
	DefaultLogLevel = "warn"
)

var initializers []func()

func init() {
	Env = &Environment{
		ClusterID:  DefaultClusterID,
		ClusterURL: DefaultClusterURL,
		Inbox:      DefaultClusterInbox,
		Outbox:     DefaultClusterOutbox,
		LogFormat:  DefaultLogFormat,
		LogLevel:   DefaultLogLevel,
	}
}

var Env *Environment

// OnInitialize sets the passed functions to be run when runtime
// is called for initialization.
func OnInitialize(y ...func()) {
	initializers = append(initializers, y...)
}

func runInitializers() {
	for _, fn := range initializers {
		fn()
	}
}

// Runtime ...
type Runtime struct {
	Run  func(*Environment)
	RunE func(*Environment) error
}

// Environment describes a runtime environment for a plugin.
// This contains information about the used NATS cluster,
// the cluster id and the topic for plugin discovery.
type Environment struct {
	ClusterID  string `mapstructure:"cluster_id"`
	ClusterURL string `mapstructure:"cluster_url"`
	Debug      bool
	Inbox      string
	LogFormat  string `mapstructure:"log_format"`
	LogLevel   string `mapstructure:"log_level"`
	Name       string
	Outbox     string
	Verbose    bool
}

// Execute ...
func (r *Runtime) Execute() error {
	runInitializers()

	if r.Run != nil {
		r.Run(Env)

		return nil
	}

	if r.RunE != nil {
		if err := r.RunE(Env); err != nil {
			return err
		}
	}

	return nil
}
