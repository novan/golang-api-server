package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/novan/golang-api-server/dic"
	"github.com/novan/golang-api-server/docs"
	"github.com/novan/golang-api-server/transport/http"
	"github.com/novan/golang-api-server/util"
)

// @title Golang API Server
// @version 1.0
// @description Golang API server example

// @contact.name Novan Adrian
// @contact.url https://novanadrian.com
// @contact.email novanadrian@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @authorizationurl JWT Authorization header using the Bearer scheme. Example: "Bearer {token}"

func main() {

	_ = util.Env("./")
	util.Log = util.NewLogger()
	util.Log.Infof("Starting service on %s environment", os.Getenv("APP_ENV"))

	// Dependency Injection
	dic.InitContainer()

	// Swagger
	docs.SwaggerInfo.Host = os.Getenv("HTTP_LISTEN")
	docs.SwaggerInfo.BasePath = "/api/v1"

	// running transport
	e := http.Run()
	r := http.NewRoutesFactory(e)
	r.Init()

	go func() {
		if err := e.Start(os.Getenv("HTTP_LISTEN")); err != nil {
			util.Log.Info("Shutting down the server...")
			util.Log.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

}
