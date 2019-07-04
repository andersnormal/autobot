<div align="center" styles="padding: 2rem;">
  <img src="https://github.com/andersnormal/autobot/blob/master/images/logo.png?raw=true" alt="Autobot"/>
</div>

# Autobot

[![Go Report Card](https://goreportcard.com/badge/github.com/andersnormal/autobot)](https://goreportcard.com/report/github.com/andersnormal/autobot)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)
[![Volkswagen](https://auchenberg.github.io/volkswagen/volkswargen_ci.svg?v=1)](https://github.com/auchenberg/volkswagen)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Autobot is your nice and friendly bot. He is here to save you from the :japanese_ogre: decepticons of #devops.

:see_no_evil: Contributions are welcome. 

## Features

* Plugable via [Pub/Sub Plugins](https://github.com/andersnormal/autobot/tree/master/pkg/plugins)

## Install

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
