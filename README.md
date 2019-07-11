<div align="center" styles="padding: 2rem;">
  <img src="https://github.com/andersnormal/autobot/blob/master/images/logo.png?raw=true" alt="Autobot"/>
</div>

# Autobot

[![License Apache 2](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Build Status](https://travis-ci.org/andersnormal/autobot.svg?branch=master)](https://travis-ci.org/andersnormal/autobot)
[![Go Report Card](https://goreportcard.com/badge/github.com/andersnormal/autobot)](https://goreportcard.com/report/github.com/andersnormal/autobot)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Autobot is your nice and friendly bot. He is here to save you from the :japanese_ogre: decepticons of #devops.

:see_no_evil: Contributions are welcome.

## Features

* Plugable architecture via [Pub/Sub Plugins](https://github.com/andersnormal/autobot/tree/master/pkg/plugins)
* Message queue for inbox /outbox via embedded [NATS Streaming](https://github.com/nats-io/stan.go)
* [Protobuf](/proto/plugin.proto) for unified messaging

## Architecture

Autobot is made of a [server](/server) and [plugins](/plugins). The server starts and embedded [Nats Streaming Server](https://github.com/nats-io/nats-streaming-server) and the provided plugins. The plugins are started with an environment that exposes two channels for publishing and subscribing to messages and some more information. The [plugins](/pkg/plugins) package exposed functions to subscribe to the `inbox` channel, which should be used to publish messages from message services and `outbox` which should publish to these services (e.g. [Slack](https://slack.com) or [Microsoft Teams](https://products.office.com/microsoft-teams/free).


## Plugins

> [godoc.org](https://godoc.org/github.com/andersnormal/autobot/pkg/plugins) for writing plugins

There are some example plugins

* [Slack](/plugins/plugin-slack/README.md)
* [Hello World](/plugins/plugin-hello-world)

## Install

### Docker

> The images are hosted on [Docker Hub](https://cloud.docker.com/u/andersnormal/repository/docker/andersnormal/autobot)

```
docker run -v $PWD/plugins:/plugins -p 8222:8222 -it andersnormal/autobot --verbose --debug --plugins /plugins
```

### Homebrew

```bash
brew install andersnormal/autobot/autobot
```

## Development

```
env GO111MODULE=on mkdir -p bin && go build -i -o bin/autobot && chmod +x bin/autobot
```

## License
[Apache 2.0](/LICENSE)
