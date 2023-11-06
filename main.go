package main

import (
	"fmt"
	"log"
	"os"

	"chatapp/app"
	"chatapp/app/client/db/mysql"
	"chatapp/app/config"
)

// variables needed
var (
	cfg    *config.Configuration
	logger *log.Logger
)

// init executed first
// start here what you need before to start the application
func init() {
	// logger for the application
	logger = getLogger("chatapp.log", "chatapp:: ")
	// configuration environment
	cfg = getConfig()
	// migrate db
	// others services should start at the biz layer instead of doing here
	logger.Printf("Starting app %s", cfg.AppName)
}

// getLogger return the logger instance
func getLogger(name, prefix string) *log.Logger {
	logpath := name
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	return log.New(file, prefix, log.LstdFlags|log.Lshortfile)
}

// getConfig return the current configuration
func getConfig() *config.Configuration {
	cfg := config.Cfg
	if cfg.Fail {
		fmt.Printf("load configuration failed: %s", cfg.FailMessage)
		os.Exit(1)
	}
	return &cfg
}

// MigrateDB the db
func MigrateDB(logger *log.Logger, cfg *config.Configuration) error {
	if cfg.DBType == "mysql" {
		mg := mysql.NewMigrtor(logger, cfg)
		if err := mg.Run(); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("migration type not found")
}

// StartApp create the application and start it
func StartApp(logger *log.Logger, cfg *config.Configuration) {
	myApp := app.New(logger)
	if err := myApp.Start(cfg); err != nil {
		if err == app.ErrorTurnOff {
			logger.Printf("Service turn off: %s", err.Error())
			return
		}
		logger.Printf("Service turn off: %s", err.Error())
		return
	}
	logger.Printf("Service turn off")
}

// main function...
// @title chat app
// @version 1.0
// @description simple api for chat in rooms
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	if cfg.Fail {
		logger.Fatalf("Error getting the configuration %s", cfg.FailMessage)
		return
	}
	// check the migration
	if cfg.MigrateDB {
		logger.Println("migrate database")
		if err := MigrateDB(logger, cfg); err != nil {
			logger.Fatalf("Error migrating the db: %+v", err)
			return
		}
	}
	// just start the application
	StartApp(logger, cfg)
}
