# jones

JONES provides with a back-end for a [jspsych](https://www.jspsych.org/) reverse correlation experience intended to compare pair of sounds or images generated with [cleese](https://github.com/neuro-team-femto/cleese). The main features are:

* enable the creation of new experiments through a few configuration files (see `examples/README.md`)
* have the jspsych (front-end) experiment interact with the server thanks to WebSockets (leveraging [golang](https://golang.org/) and [gorilla/websocket](https://github.com/gorilla/websocket))
* deployment made easy (1): build project, transfer the binary and a few folders
* deployment made easy (2): no database needed (all configuration and state being saved to text files)
* deployment made easy (3): the `jones` binary comes with a HTTP server for static files (JS, CSS)

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

2. Build the `jones` binary:

```sh
# creates the 'jones' binary
go build
```

Check [Build options](#build-options) to build for different platforms.

3. Build front-end assets (to the `public` folder) thanks to the `jones` binary:

```sh
# with an environment variable:
APP_MODE=BUILD_FRONT ./jones
# or with a command line argument:
./jones -APP_MODE BUILD_FRONT
```

4. Transfer to server the `jones` binary, and the `data`, `examples` and `public` folders:

```
jones      -> jspsych WebSocket back-end to manage experiments + HTTP server for front-end assets 
data/      -> contains live experiments data (configuration and results)
examples/  -> (optional) contains example experiment configurations that you may copy/paste then edit
public/    -> js/css assets served by HTTP server including those created with APP_MODE=BUILD_FRONT
```

Other files and folders in this project are only needed for the build (steps 2 and 3) and don't need to be transferred to the server.

Please note the `jones` binary will automatically create an additional `state` folder to manage internal state.

4. Run `jones` (with a user with write permissions on local folder and below), at least specifying from what origins WebSockets connections are allowed:

```sh
APP_ORIGINS=https://example.com ./jones
```

Check other available settings in the [Environment variables](#environment-variables) section.

## Environment variables

* `APP_PORT=9000` (defaults to 8100) to set port listen by `jones` server
* `APP_ORIGINS=https://example.com` to declare a comma separated list of allowed origins for WebSocket connections (`http://localhost:8100` and `https://localhost:8100` are allowed by default if `APP_ORIGINS` is not set)
* `APP_WEB_PREFIX=/path` (empty by default) needed if, depending on your server configuration, `jones` is served under a given path, for instance `https://example.com/path`
* `APP_MODE=DEV` to enable development mode (watch JS files to trigger builds, enhanced logs and allow `http://localhost:8100` and `https://localhost:8100` origins)
* `APP_MODE=BUILD_FRONT` builds front-end assets but do not start server

## Command line arguments

It's also possible to run the projet with the `APP_MODE` command line argument, in that case it will have priority over the corresponding environment variable if both are defined, for instance:

```sh
go run main.go -APP_MODE BUILD_FRONT
```

## Create a new experiment

Create a new folder in `data/` and follow the instructions in `examples/README.md`.

## Proxy setup

You may prefer to serve `jones` behind a HTTP proxy, for instance proxying a given public sub/domain to the declared `APP_PORT` of the locally running `jones`.

In that case you need to allow for WebSockets upgrade and to fine-tune the default WebSockets timeout (taking into consideration possible idle periods of participants) at the proxy side.

For instance if you use nginx, you may consider the following directives:

```
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection $connection_upgrade;
# increase default timeout
proxy_send_timeout 10m;
proxy_read_timeout 10m;
```

## Process control

You may manage `jones` execution with [Supervisor](http://supervisord.org/) or [pm2](https://pm2.keymetrics.io/docs/usage/quick-start/) (for instance for auto-restarts). Here is a configuration example for Supervisor:

```
[program:jones]
directory=/home/deploy/jones
command=/home/deploy/jones/jones
stdout_logfile=/home/deploy/jones/out.log
stderr_logfile=/home/deploy/jones/err.log
environment=APP_WEB_PREFIX=/jones,APP_ORIGINS=https://example.com
user=deploy
autostart=true
autorestart=true
```

## Build options

Check available options if you build for a different machine (see [more](https://golang.org/doc/install/source#environment)) for instance:

```sh
GOOS=linux GOARCH=amd64 go build
```