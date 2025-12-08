package main

import (
	"flag"
	"log"
	"os"

	server "github.com/zercle/gofiber-skelton/internal/infrastructure"
	"github.com/zercle/gofiber-skelton/pkg/config"
)

// @title Waste Management Service API
// @version 1.0
// @description This is the API documentation for the Waste Management Service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var (
	version string
	build   string
	runEnv  string
)

func init() {
	// read running flag
	Testttst
	if len(os.Getenv("ENV")) != 0 {
		runEnv = os.Getenv("ENV")
	} else {
		flagEnv := flag.String("env", "dev", "A config file name without .env")
		flag.Parse()
		runEnv = *flagEnv
	}
	// load config by running flag
	if err := config.LoadConfig(runEnv); err != nil {
		log.Fatalf("error while loading the env:\n %+v", err)
	}
}

func main() {

	// init server
	server, err := server.NewServer(version, build, runEnv)
	if err != nil {
		log.Fatalf("error while create server:\n %+v", err)
	}

	server.Run()
}
