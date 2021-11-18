# revcor

CAUTION: deploy instructions currently under work

This project provides with a back-end for a [jspsych](https://www.jspsych.org/) revcor experience. The main features are:

* enable the creation of new experiments through a few configuration files (see `data/README.md`)
* interact with the experiment thanks to WebSockets (leveraging [golang](https://golang.org/) and [gorilla mux](https://github.com/gorilla/mux))
* deployment made easy (1): build project, transfer the binary and a few configuration files
* deployment made easy (2): no database needed (all configuration and state being saved to text files)
* deployment made easy (3): the `revcor` binary comes with a HTTP server for static files (JS, CSS)

## Build and deploy

1. Clone the current project

2. Build it:

```
go build
```

Check available options if you build for a different machine (see [more](https://golang.org/doc/install/source#environment)) for instance:

```
GOOS=linux GOARCH=amd64 go build
```

3. Transfer the `revcor` binary (created by step 2) to your server along with the following folders: `data` and `front/public` to obtain the following layout:

```
revcor
data/         -> contains experiments configuration and generated data
front/public  -> js/css assets served by HTTP server
```

4. Launch `revcor` at least specifying from what origins WebSockets connections are allowed:

```
APP_ORIGINS=https://example.com ./revcor
```

Check other available settings in the [Environment Variables](#environment-variables) section.

## Environment variables

* `APP_PORT=9000` (defaults to 8100) to set port listen by `revcor` server
* `APP_ORIGINS=https://example.com` to declare a comma separated list of allowed origins for WebSocket connections
* `APP_WEB_PREFIX=/path` if, depending on your server configuration, `revcor` is served under a given path, for instance `https://example.com/path`
* `APP_ENV=DEV` to enable development mode (set a few allowed origins on localhost, watch JS files to trigger builds, enhance logs)
* `APP_LOG_FILE=revcor.log` (defaults to none) to declare a file to write logs to (fails silently if file can't be opened)
* `APP_LOG_STDOUT=true` (defaults to false) to print logs to Stdout (if `APP_LOG_FILE` is also set, logs are written to both)

## Additional setup

You may prefer to serve `revcor` behind a HTTP proxy, for instance as an nginx server block that forwards to the declared `APP_PORT`.

You may also manage `revcor` execution with supervisord or pm2.