module github.com/andersnormal/autobot

go 1.12

require (
	github.com/andersnormal/pkg v0.0.0-20190521194814-d257d6a24e99
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/docker/libkv v0.2.1
	github.com/golang/protobuf v1.3.1
	github.com/google/uuid v1.1.1
	github.com/ianschenck/envflag v0.0.0-20140720210342-9111d830d133
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/lusis/go-slackbot v0.0.0-20180109053408-401027ccfef5 // indirect
	github.com/lusis/slack-test v0.0.0-20190426140909-c40012f20018 // indirect
	github.com/nats-io/gnatsd v1.4.1
	github.com/nats-io/go-nats v1.7.2 // indirect
	github.com/nats-io/nats-server v1.4.1 // indirect
	github.com/nats-io/nats-streaming-server v0.15.1
	github.com/nats-io/nats.go v1.8.1
	github.com/nats-io/stan.go v0.4.5
	github.com/nlopes/slack v0.5.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	google.golang.org/grpc v1.21.0
)

replace github.com/golang/protobuf v1.3.1 => github.com/golang/protobuf v1.2.0
