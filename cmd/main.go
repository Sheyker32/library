package main

import (
	_ "library/cmd/docs"
	"library/config"
	"library/postgres"
	"library/run"
	"os"

	"library/logging"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// @title						Swagger Petstore
// @version					1.0
// @description				Documentation for petstore
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	logger := logging.GetLogger()

	conf, err := config.LoadDBConfig()
	if err != nil {
		logger.Fatal("Failed to load DB config: ", zap.Error(err))
	}

	db := postgres.NewPostgresDB(conf, logger)
	defer db.Close()

	app := run.NewApp(db, logger)

	exitCode := app.
		Bootstrap().
		Run()
	os.Exit(exitCode)
}
