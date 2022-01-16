package main

import (
	stdlog "log"
	"os"

	"github.com/joho/godotenv"
	"github.com/josephspurrier/ambient"
	"github.com/josephspurrier/ambient/cmd/myapp/app"
	"github.com/josephspurrier/ambient/lib/aesdata"
	"github.com/josephspurrier/ambient/lib/cloudstorage"
	"github.com/josephspurrier/ambient/lib/envdetect"
	"github.com/josephspurrier/ambient/plugin/logger/zaplogger"
)

var (
	appName    = "myapp"
	appVersion = "1.0"
)

func init() {
	// Verbose logging with file name and line number for the standard logger.
	stdlog.SetFlags(stdlog.Lshortfile)
}

func main() {
	// Load the .env file if running in local dev mode.
	if envdetect.RunningLocalDev() {
		err := godotenv.Load()
		if err != nil {
			stdlog.Fatalf("app: error loading .env file: %v\n", err.Error())
		}
	}

	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		stdlog.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	// Determine cloud storage engine for site and session information.
	storage := cloudstorage.StorageBasedOnCloud(app.StorageSitePath,
		app.StorageSessionPath)

	// Create the ambient app.
	plugins := app.Plugins()
	ambientApp, log, err := ambient.NewApp(appName, appVersion,
		zaplogger.New(),
		ambient.StoragePluginGroup{
			Storage:    storage,
			Encryption: aesdata.NewEncryptedStorage(secretKey),
		},
		plugins)
	if err != nil {
		if log != nil {
			// Use the logger if it's available.
			log.Fatal("", err.Error())
		} else {
			// Else use the standard logger.
			stdlog.Fatalln(err.Error())
		}
	}

	// Set the log level.
	// ambientApp.SetLogLevel(ambient.LogLevelDebug)

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		log.Fatal("", err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe(mux)
}
