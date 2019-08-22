<div align="center" styles="padding: 2rem;">
  <img src="https://github.com/andersnormal/autobot/blob/master/images/logo.png?raw=true" alt="Autobot"/>
</div>

# Autobot

[![GoDoc](https://godoc.org/github.com/narqo/go-badge?status.svg)](https://godoc.org/github.com/andersnormal/autobot)
[![License Apache 2](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Build Status](https://travis-ci.org/andersnormal/autobot.svg?branch=master)](https://travis-ci.org/andersnormal/autobot)
[![Go Report Card](https://goreportcard.com/badge/github.com/andersnormal/autobot)](https://goreportcard.com/report/github.com/andersnormal/autobot)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

https://godoc.org/github.com/andersnormal/autobot

Autobot is your nice and friendly bot. He is here to save you from the :japanese_ogre: decepticons of #devops.

:see_no_evil: Contributions are welcome.

## Features

* Plugable architecture via [Pub/Sub Plugins](https://github.com/andersnormal/autobot/tree/master/pkg/plugins)
* Message queue for inbox /outbox via embedded [NATS Streaming](https://github.com/nats-io/stan.go)
* [Protobuf](/proto/plugin.proto) for unified messaging

## Architecture

Autobot is made of a [server](/server) and [plugins](/plugins). The server starts and embedded [Nats Streaming Server](https://github.com/nats-io/nats-streaming-server) and the provided plugins. The plugins are started with an environment that exposes two channels for publishing and subscribing to messages and some more information. The [plugins](/pkg/plugins) package exposed functions to subscribe to the `inbox` channel, which should be used to publish messages from message services and `outbox` which should publish to these services (e.g. [Slack](https://slack.com) or [Microsoft Teams](https://products.office.com/microsoft-teams/free).

It is also built up on a 3 factor architecture. Which means that is has a global state that is changed by the messages flowing in the system. All authentication and authorization is handeled via the global state.

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

### Docker Compose

Because Autobot is using pub/sub to communicate with its plugins they can be run independently in their own containers. [Anders Normal](https://cloud.docker.com/u/andersnormal) contains the plugins in containers. [docker-compose.yml](/docker-compose.yml) contains an example to run the provided plugins in their containers.

> you should change [.env](/.env) for your specific setup

```
# setart the containers
docker-compose up
```

## Development

> we use [Picasso](https://github.com/andersnormal/picasso) for build automation 

You can build the [Protobuf](/proto) by running 

```
picasso proto
```

The server is build by running

```
picasso build/server
```

The options of the server can be shown by `./server --help`.

## License
[Apache 2.0](/LICENSE)
