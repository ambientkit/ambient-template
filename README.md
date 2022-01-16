# Sample App using Ambient

This repository contains a sample app to demonstrate how to use the [Ambient](https://github.com/josephspurrier/ambient) pluggable web framework.

### What is it?

Ambient is a framework in Go for building web apps using plugins. You can use the plugins already included to stand up a blog just like the [Bear Blog](https://bearblog.dev/) or create your own plugins to build your own web app. Plugins can be enabled/disabled while the app is running which means routes as well as middleware can also modified without restarting the app. Plugins must be granted permissions above being enabled which provides you with better control over your web app.

You can read more why the framework was created [here](https://github.com/josephspurrier/ambient).

Use the [Deployment Guide](DEPLOYMENT.md) to deploy serverless on Google Cloud (Cloud Run), AWS (App Runner), or Azure (Functions).

Use the [Plugin Development Guide](PLUGIN.md) to build your own plugins.

## Quickstart on Local

To test out the sample web app:

- Clone the repository: `git clone git@github.com:josephspurrier/ambient-template.git`
- Create a new file called `.env` in the root of the repository with this content:

```bash
# App version.
AMB_APP_VERSION=1.0
# Set this to any value to allow you to do testing locally without cloud access.
AMB_LOCAL=true
# Optional: Enable the Dev Console that amb connects to. Default is: false
AMB_DEVCONSOLE_ENABLE=true
# Optional: Set the URL for the Dev Console that amb connects to. Default is: http://localhost
# AMB_DEVCONSOLE_URL=http://localhost
# Optional: Set the port for the Dev Console that amb connects to. Default is: 8081
# AMB_DEVCONSOLE_PORT=8081
# Session key to encrypt the cookie store. Generate with: make privatekey
AMB_SESSION_KEY=
# Password hash that is base64 encoded. Generate with: make passhash passwordhere
AMB_PASSWORD_HASH=

# Optional: set the web server port.
# PORT=8080
# Optional: set the time zone from here:
# https://golang.org/src/time/zoneinfo_abbrs_windows.go
# AMB_TIMEZONE=America/New_York
# Optional: set the URL prefix if behind a proxy.
# AMB_URL_PREFIX=/api
```

- To create the session and site files in the storage folder, run: `make storage`
- To start the webserver on port 8080, run: `make run-env`

The login page is located at: http://localhost:8080/login.

To login, you'll need:

- the default username is: `admin`
- the password from the .env file for which the `AMB_PASSWORD_HASH` was derived

Once you are logged in, you should see a new menu option call `Plugins`. From this screen, you'll be able to use the Plugin Manager to make changes to state, permissions, and settings for all plugins.

### Local Development Flags

You can set the web server `PORT` to values other than `8080`.

When `AMB_LOCAL` is set to `true`:

- data storage will use the local filesystem instead of Google Cloud Storage
- if you try to access the app, it will listen on all IPs/addresses, instead of redirecting like it does in production

You can use `envdetect.RunningLocalDev()` to detect if the flag is set to true or not.

When `AMB_TIMEZONE` is set to a timezone like `America/New_York`, the app will use that timezone. This is required if using time-based packages like MFA.

When `AMB_URL_PREFIX` is set to a path like `/api`, the app will serve requests from `/api/...`. This is helpful if you are running behind a proxy or are hosting multiple websites from a single URL.

### App Settings

In the main.go file, you can modify your log level with `SetLogLevel()`:

```go
ambientApp, err := ambient.NewApp(...)
ambientApp.SetLogLevel(ambient.LogLevelDebug)
ambientApp.SetLogLevel(ambient.LogLevelInfo)
ambientApp.SetLogLevel(ambient.LogLevelError)
ambientApp.SetLogLevel(ambient.LogLevelFatal)
```

You can enable `span` tags around HTML elements to determine which content is loaded from which plugins with `SetDebugTemplates()`:

```go
ambientApp, err := ambient.NewApp(...)
ambientApp.SetDebugTemplates(true)
```

You can disable template escaping with `SetEscapeTemplates()`:

```go
ambientApp, err := ambient.NewApp(...)
ambientApp.SetEscapeTemplates(false)
```

## Development Workflow

If you would like to make changes to the code with hot reloading capabilities, I recommend [`air`](https://github.com/cosmtrek/air) to help streamline your workflow.

```bash
# Install air to allow hot reloading so you can make changes quickly.
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
```

You can then use this command to start the web server and monitor for changes:

```bash
# Start hot reload. The web app should be available at: http://localhost:8080
air
```