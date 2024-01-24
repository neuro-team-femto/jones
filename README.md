# revcor

This project provides with a back-end for a [jspsych](https://www.jspsych.org/) revcor experience intended to compare pair of sounds or images generated with [cleese](https://github.com/neuro-team-femto/cleese). The main features are:

* enable the creation of new experiments through a few configuration files (see `examples/README.md`)
* have the jspsych (front-end) experiment interact with the server thanks to WebSockets (leveraging [golang](https://golang.org/) and [gorilla/websocket](https://github.com/gorilla/websocket))
* deployment made easy (1): build project, transfer the binary and a few folders
* deployment made easy (2): no database needed (all configuration and state being saved to text files)
* deployment made easy (3): the `revcor` binary comes with a HTTP server for static files (JS, CSS)

## Run in development

Clone the repository then:

```sh
# with an environment variable:
APP_MODE=DEV go run main.go
# or with a command line argument:
go run main.go -APP_MODE DEV
```

Then go to http://localhost:8100/xp/example/new

## Build and deploy

1. Clone the repository

2. Build the `revcor` binary:

```sh
# creates the 'revcor' binary
go build
```

Check [Build options](#build-options) to build for different platforms.

3. Build front-end assets (to the `public` folder) thanks to the `revcor` binary:

```sh
# with an environment variable:
APP_MODE=BUILD_FRONT ./revcor
# or with a command line argument:
./revcor -APP_MODE BUILD_FRONT
```

4. Transfer to server the `revcor` binary, and the `data`, `examples` and `public` folders:

```
revcor     -> jspsych WebSocket back-end to manage experiments + HTTP server for front-end assets 
data/      -> contains live experiments data (configuration and results)
examples/  -> (optional) contains example experiment configurations that you may copy/paste then edit
public/    -> js/css assets served by HTTP server including those created with APP_MODE=BUILD_FRONT
```

Other files and folders in this project are only needed for the build (steps 2 and 3) and don't need to be transferred to the server.

Please note the `revcor` binary will automatically create an additional `state` folder to manage internal state.

4. Run `revcor` (with a user with write permissions on local folder and below), at least specifying from what origins WebSockets connections are allowed:

```sh
APP_ORIGINS=https://example.com ./revcor
```

Check other available settings in the [Environment variables](#environment-variables) section.

## Environment variables

* `APP_PORT=9000` (defaults to 8100) to set port listen by `revcor` server
* `APP_ORIGINS=https://example.com` to declare a comma separated list of allowed origins for WebSocket connections (`http://localhost:8100` and `https://localhost:8100` are allowed by default if `APP_ORIGINS` is not set)
* `APP_WEB_PREFIX=/path` (empty by default) needed if, depending on your server configuration, `revcor` is served under a given path, for instance `https://example.com/path`
* `APP_MODE=DEV` to enable development mode (watch JS files to trigger builds, enhanced logs and allow `http://localhost:8100` and `https://localhost:8100` origins)
* `APP_MODE=BUILD_FRONT` builds front-end assets but do not start server

## Command line arguments

It's also possible to run the projet with the `APP_MODE` command line argument, in that case it will have priority over the corresponding environment variable if both are defined, for instance:

```sh
go run main.go -APP_MODE BUILD_FRONT
```

## Create a new experiment

Create a new folder in `data/` and follow the instructions in `examples/README.md`.

## Additional setup

You may prefer to serve `revcor` behind a HTTP proxy, for instance as an nginx server block that forwards to the declared `APP_PORT`.

You may also manage `revcor` execution with [Supervisor](http://supervisord.org/) or [pm2](https://pm2.keymetrics.io/docs/usage/quick-start/). Here is a configuration example for Supervisor:

```
[program:revcor]
directory=/home/deploy/revcor
command=/home/deploy/revcor/revcor
stdout_logfile=/home/deploy/revcor/out.log
stderr_logfile=/home/deploy/revcor/err.log
environment=APP_WEB_PREFIX=/revcor,APP_ORIGINS=https://example.com
user=deploy
autostart=true
autorestart=true
```

## Build options

Check available options if you build for a different machine (see [more](https://golang.org/doc/install/source#environment)) for instance:

```sh
GOOS=linux GOARCH=amd64 go build
```