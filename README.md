<div align="center" styles="padding: 2rem;">
  <img src="https://github.com/andersnormal/autobot/blob/master/images/logo.png?raw=true" alt="Autobot"/>
</div>

# Autobot

<a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-42%25-brightgreen.svg?longCache=true&style=flat)</a>
[![Build Status](https://travis-ci.org/andersnormal/autobot.svg?branch=master)](https://travis-ci.org/andersnormal/autobot)
[![Go Report Card](https://goreportcard.com/badge/github.com/andersnormal/autobot)](https://goreportcard.com/report/github.com/andersnormal/autobot)
[![License Apache 2](https://img.shields.io/badge/License-Apache2-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)

Autobot is your nice and friendly bot. He is here to save you from the :japanese_ogre: decepticons of #devops and other evil.

:see_no_evil: Contributions are welcome.

## Features

* Plugable architecture via [Pub/Sub Plugins](https://github.com/andersnormal/autobot/tree/master/pkg/plugins)
* Message queue for the bot inbox /outbox via embedded [NATS Streaming](https://github.com/nats-io/stan.go)
* [Protobuf](/proto/plugin.proto) for unified messages
* Plugins (e.g. Slack) but many more that you can build

## Architecture

Autobot is made of a [server](/server) and [plugins](/plugins). The server starts an embedded [Nats Streaming Server](https://github.com/nats-io/nats-streaming-server). The plugins subscribe and publish message to the provided queues. They can be run in their custom Docker containers. The plugins are started with an environment that exposes two channels for publishing and subscribing to messages and some more information. The [plugins](/pkg/plugins) package exposed functions to subscribe to the `inbox` channel, which should be used to publish messages from message services and `outbox` which should publish to these services (e.g. [Slack](https://slack.com) or [Microsoft Teams](https://products.office.com/microsoft-teams/free)).

<pre>
    Slack /                                                                
 Microsoft Teams               +-----------------------------------+    
        ^                    +-|---------------------------------+ |    
        |                    | |            Server               | |    
        |                    | | +---------------+               | |    
 +------|---------\          | | |     NATS      |               | |    
 |   Plugin    |   -----\    | | |               |               | |    
 +-------------<-\       -----\| |+-------------+|               | |    
                  ----\    ------->   Inbox     ||               | |    
                   -------\  | | |+-------------+|               | |    
 +-------------<--/        ----\ |+-------------+|               | |    
 |   Plugin    |         --------->   Outbox    ||               | |    
 +----------------------/    | | |+-------------+|               | |    
                             | | +---------------+               | |    
                             | +-----------------------------------+    
                             +-----------------------------------+  
</pre>

### Clustering

NATS Streaming Server supports [clustering](https://nats-io.github.io/docs/nats_streaming/clustering/clustering.html) and data replication. Because Autobot embedds the streaming server it supports clustering as the mechanism for high availability. Autobot supports the two modes of clustering. We recommend the "auto" mode in which you specify the list of peers to each started node. They elect a leader and start the replication. Because it uses the [RAFT consensus algorithm](https://raft.github.io/) it needs an uneven number of peers in the cluster. 3 or 5 are sufficient.

```bash
# this create a cluster consiting of 3 nodes
docker-compose up
```

## Plugins

> [godoc.org](https://godoc.org/github.com/andersnormal/autobot/pkg/plugins) for writing plugins

There are some example plugins

* [Slack](/plugins/plugin-slack/README.md)
* [Hello World](/plugins/plugin-hello-world)

Plugins are either run and managed by the controller or they are run as individual processes in a container and are individually managed. The `--plugins` flag specifies directories or individual plugins to be run and managed by the controller. 

Plugins can be either configured by the automatically exposed command line parameters or the prefixed environment variables.

Example for the [Slack Plugin](https://github.com/andersnormal/autobot/plugins/plugin-slack/README.md):

```bash
SLACK_TOKEN=
AUTOBOT_CLUSTER_ID=autobot
AUTOBOT_CLUSTER_URL=nats://controller:4222
AUTOBOT_LOG_FORMAT=json
AUTOBOT_LOG_LEVEL=info
AUTOBOT_DEBUG=true
AUTOBOT_VERBOSE=true
```

There are two log formats supported `text` (default) and `json`. The log levels reflect the [logrus](https://github.com/sirupsen/logrus/blob/4f5fd631f16452fbd023813c1eb7dbd67130cb0c/logrus.go#L93) levels.

## Example

> The images are hosted on [Docker Hub](https://cloud.docker.com/u/andersnormal/repository/docker/andersnormal/autobot)
> you should change [.env](/.env) for your specific setup

This example uses [Docker Compose](https://docs.docker.com/compose/) and you will need a [Slack Bot](https://api.slack.com/bot-users) Token (e.g. xob-xxxx).

You should provide this token in the [.env](https://github.com/andersnormal/autobot/.env) file. Which is used to configure the plugins containers. Because Autobot is using pub/sub to communicate with its plugins they can be run independently in their own containers. [Anders Normal](https://cloud.docker.com/u/andersnormal) contains the plugins in containers. 

```bash
# start the containers
docker-compose up
```

You should now see your Slack Bot connect in your Workspace and can send him a direct message which he will respond to with `hello world`.

## Development

> we use [Picasso](https://github.com/andersnormal/picasso) for build automation 

You can build the [Protobuf](/proto) by running 

```bash
picasso proto
```

We use a specific version of proto package generator. In order to build it with this version you will have to install it as follows

```bash
GIT_TAG="v1.2.0" # change as needed
go get -d -u github.com/golang/protobuf/protoc-gen-go
git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG
go install github.com/golang/protobuf/protoc-gen-go
```

The server is build by running

```bash
picasso build/server
```

The options of the server can be shown by `./server --help`.

## License

[Apache 2.0](/LICENSE)
