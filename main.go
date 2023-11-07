package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"chatapp/app"
	"chatapp/app/biz"
	"chatapp/app/client/db/mysql"
	"chatapp/app/client/rabbitmq"
	"chatapp/app/command"
	"chatapp/app/config"
	"chatapp/pkg/utils"
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
	logger = utils.GetLogger("chatapp.log", "chatapp:: ")
	// configuration environment
	cfg = utils.GetConfig()
	// migrate db
	// others services should start at the biz layer instead of doing here
	logger.Printf("Starting app %s", cfg.AppName)
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

// StartBot ...
func StartBot(logger *log.Logger, cfg *config.Configuration) {
	logger = utils.GetLogger("botcommand", "cmd :: ")
	biz := biz.New(logger, nil)
	logger.Println("start bot: ", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host)
	rmq := rabbitmq.NewRabbit(logger, cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host)
	biz.SetRabbit(rmq)
	ncmd := command.NewCommand(logger, cfg, biz, rmq)
	ncmd.Start()
	logger.Printf("running the command api")
	ncmd.Run()
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

	operation := flag.String("t", "api", "type: 'api' to start the http.rest web \n 'bot' to start the bot ")
	flag.Parse()
	fmt.Println("operation: ", *operation)
	switch strings.ToLower(*operation) {
	case "api":
		// just start the application
		StartApp(logger, cfg)
	case "bot":
		// start the bot
		StartBot(logger, cfg)
	}
}
