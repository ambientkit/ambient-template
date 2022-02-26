// Package main App API
//
// This Swagger specification describes the application API.
//
// Swagger 2.0 Spec - generated by [go-swagger](https://github.com/go-swagger/go-swagger)
//
// Schemes: http, https
// Host: localhost:8080
// BasePath: /
// Version: VERSIONPLACEHOLDER
// License: Copyright © 2022 Company. All rights reserved.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// SecurityDefinitions:
// token:
//   type: apiKey
//   name: Authorization
//   in: header
//   description: "Enter your API key prefixed by the word \"Bearer\"."
//
// swagger:meta
package main

import (
	"log"
	"os"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient-template/cmd/myapp/app"
	"github.com/ambientkit/ambient/pkg/ambientapp"
	"github.com/ambientkit/ambient/pkg/envdetect"
	"github.com/ambientkit/plugin/logger/zaplogger"
	"github.com/ambientkit/plugin/pkg/aesdata"
	"github.com/ambientkit/plugin/pkg/cloudstorage"
	"github.com/joho/godotenv"
)

var (
	appName    = "myapp"
	appVersion = "1.0"
)

func init() {
	// Verbose logging with file name and line number for the standard logger.
	log.SetFlags(log.Lshortfile)
}

func main() {
	// Load the .env file if set to load the file.
	if envdetect.LoadDotEnv() {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("app: error loading .env file: %v\n", err.Error())
		}
	}

	// Get the environment variables.
	secretKey := os.Getenv("AMB_SESSION_KEY")
	if len(secretKey) == 0 {
		log.Fatalf("app: environment variable missing: %v\n", "AMB_SESSION_KEY")
	}

	// Determine cloud storage engine for site and session information.
	storage := cloudstorage.StorageBasedOnCloud(app.StorageSitePath,
		app.StorageSessionPath)

	// Create the ambient app.
	plugins := app.Plugins()
	ambientApp, logger, err := ambientapp.NewApp(appName, appVersion,
		zaplogger.New(),
		ambient.StoragePluginGroup{
			Storage:    storage,
			Encryption: aesdata.NewEncryptedStorage(secretKey),
		},
		plugins)
	if err != nil {
		// Use the standard logger.
		log.Fatalln(err.Error())
	}

	// Set the log level.
	// ambientApp.SetLogLevel(ambient.LogLevelDebug)

	// Load the plugins and return the handler.
	mux, err := ambientApp.Handler()
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Start the web listener.
	ambientApp.ListenAndServe(mux)
}
