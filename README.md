# revcor

This project provides with a back-end for a [jspsych](https://www.jspsych.org/) revcor experience intended to compare pair of sounds generated with [clesse](https://github.com/creamlab/cleese). The main features are:

* enable the creation of new experiments through a few configuration files (see `data/README.md`)
* interact with the experiment thanks to WebSockets (leveraging [golang](https://golang.org/) and [gorilla mux](https://github.com/gorilla/mux))
* deployment made easy (1): build project, transfer the binary and a few folders
* deployment made easy (2): no database needed (all configuration and state being saved to text files)
* deployment made easy (3): the `revcor` binary comes with a HTTP server for static files (JS, CSS)

## Build and deploy

1. Clone the current project

2. Build the `revcor` binary:

```
go build
```

Check [Build options](#build-options) to build for different platforms.

3. Build front-end assets (to the `public` folder) thanks to the `revcor` binary:

```
APP_ENV=BUILD_FRONT ./revcor
```

4. Transfer to server the `revcor` binary, and the whole `data` and `public` folders:

```
revcor     -> HTTP server for front-end assets + jspsych WebSocket back-end to manage experiments 
data/      -> contains experiments configuration and generated data
public/    -> js/css assets served by HTTP server
```

Other files and folders in this project are only needed for the build (steps 2 and 3) and thus don't have to be transferred to the server.

Please note the `revcor` binary will automatically create an additional `state` folder to manage internal state.

4. Run `revcor`, at least specifying from what origins WebSockets connections are allowed:

```
APP_ORIGINS=https://example.com ./revcor
```

Check other available settings in the [Environment variables](#environment-variables) section.

## Environment variables

* `APP_PORT=9000` (defaults to 8100) to set port listen by `revcor` server
* `APP_ORIGINS=https://example.com` to declare a comma separated list of allowed origins for WebSocket connections
* `APP_WEB_PREFIX=/path` if, depending on your server configuration, `revcor` is served under a given path, for instance `https://example.com/path`
* `APP_ENV=DEV` to enable development mode (set a few allowed origins on localhost, watch JS files to trigger builds, enhance logs)
* `APP_ENV=BUILD_FRONT` builds front-end assets but do not start server
* `APP_ADMIN_LOGIN` and `APP_ADMIN_PASSWORD` credentials to access `/admin` pages

## Create a new experiment

Create a new folder in `data` and follow the instructions in `data/README.md`.

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

```
GOOS=linux GOARCH=amd64 go build
```