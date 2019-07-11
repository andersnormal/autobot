<div align="center" styles="padding: 2rem;">
  <img src="https://github.com/andersnormal/autobot/blob/master/images/logo.png?raw=true" alt="Autobot"/>
</div>

# Autobot Slack Plugin

[![Documentation](https://godoc.org/github.com/andersnormal/autobot?status.svg)](https://godoc.org/github.com/andersnormal/autobot/plugins/plugin-slack)

This is a plugin for :robot: [Autobot](https://github.com/andersnormal/autobot) connecting a [Slack Bot](https://slack.com) to it.

:see_no_evil: Contributions are welcome.

## Usage

The plugin uses the `SLACK_TOKEN` environment variable to connect to the [Slack RTM API](https://api.slack.com/rtm). The plugin can be run with out the server, but it is recommended to do so.

> if you want to run the plugin without the server, you have to configure the environment variables [godoc.org](https://godoc.org/github.com/andersnormal/autobot/pkg/plugins) which are set by the server

```
# Example of running the plugin
./server --verbose --debug --env SLACK_TOKEN=your_token --plugins ../plugins/plugin-slack
```

> when you run the plugin via the server, all `--env` parameters are passed to the started plugins.

## Build

```
go build
```

## License
[Apache 2.0](/LICENSE)
