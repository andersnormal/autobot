# :art: Autobot

Autobot is here to save you from the :japanese_ogre: decepticons of #devops.

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
