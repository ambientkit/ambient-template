# Sample App using Ambient

Thanks for visiting! All docs are available [here](https://ambientkit.github.io/docs/).

This repository contains a sample app to demonstrate how to use the [Ambient](https://github.com/ambientkit/ambient) pluggable web framework. Ambient is a framework in Go for building web apps using plugins. You can use the plugins already included to stand up a blog just like the [Bear Blog](https://bearblog.dev/) or create your own plugins to build your own web app. Plugins can be enabled/disabled while the app is running which means routes as well as middleware can also modified without restarting the app. Plugins must be granted permissions above being enabled which provides you with better control over your web app. You can read why the Ambient framework was created [here](https://ambientkit.github.io/docs/docs/faq).

Use the [Deployment Guide](DEPLOYMENT.md) to deploy serverless on Google Cloud (Cloud Run), AWS (App Runner), or Azure (Functions).

Use the [Plugin Development Guide](https://ambientkit.github.io/docs/docs/plugins/overview) to build your own plugins.

## Quickstart on Local

To test out the sample web app, you can run these commands below. You can also view screenshots of the app [here](https://ambientkit.github.io/docs/docs/introduction/screenshots).

```bash
# Build the Ambient interactive CLI (amb) in the current folder.
bash -c "$(curl -fsSL https://raw.githubusercontent.com/ambientkit/amb/main/bash/install.sh)"

# Run the app.
./amb

# Clone the ambient template by typing this command and pressing Enter.
createapp

# Exit by typing `exit` or pressing Ctrl+D.
exit

# Change to the new project folder.
cd ambapp

# Download the Go dependencies.
go mod download

# Create the session and site files in the storage folder.
make storage

# Create the .env file.
make env

# Generate a new private key.
make privatekey >> .env

# Generate a new password hash (replace with your password).
make passhash passwordhere >> .env

# Start the webserver on port 8080 (local development with no Docker).
make run-env
```

The login page is located at: http://localhost:8080/login.

To login, you'll need:

- the default username is: `admin`
- the password from the .env file for which the `AMB_PASSWORD_HASH` was derived

Once you are logged in, you should see a new menu option call `Plugins`. From this screen, you'll be able to use the Plugin Manager to make changes to state, permissions, and settings for all plugins.

## Container Deployment

To test out the sample app in Docker, you can run these commands:

```bash
# Build the Docker container.
make build

# Test running the Docker container.
make run
```

The login page is located at: http://localhost:8080/login.

## Automatic Rebuilds

If you would like to make changes to the code that rebuilds automatically, it's recommended to use [`air`](https://github.com/cosmtrek/air) to help streamline your workflow.

```bash
# Install air to project bin folder.
make air-install

# Start air to monitor code changes. The web app should be available at:
# http://localhost:8080
air
```

## Swagger Generation

You can easily [generate a Swagger spec](https://goswagger.io/use/spec.html) for the API from annotations in the code:

```bash
# Install Swagger to project bin folder.
make swagger-install

# Generate Swagger spec: swagger.json.
make swagger

# Serve the Swagger spec and open a browser window to view and make requests.
# You need to enable the `healthcheck` and `cors` plugins for this testable UI
# to function properly. You will also need the app running.
make swagger-serve
```

You can also generate SDKs and a CLI to interact with the API. You can read more about it [here](https://goswagger.io/generate/requirements.html).

## Environment Variables

You can read about environment variables [here](https://ambientkit.github.io/docs/docs/architecture/envars).

## App Settings

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
