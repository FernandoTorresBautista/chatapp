package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"chatapp/app/api"
	"chatapp/app/biz"
	"chatapp/app/client/db"
	"chatapp/app/client/db/mysql"
	"chatapp/app/client/rabbitmq"
	"chatapp/app/config"
)

// Api structure
type Api struct {
	logger     *log.Logger
	Ctx        context.Context
	Cancel     context.CancelFunc
	bizLayer   *biz.Biz
	repository db.Repository
	apis       api.ApiHandle
}

// ErrorTurnOff ok
var ErrorTurnOff = fmt.Errorf("application turnoff ok")

// New app instance
func New(logger *log.Logger) *Api {
	return &Api{
		logger: logger,
	}
}

func (a *Api) setRepo(cfg *config.Configuration) {
	// for the moment this api only support mysql...
	switch cfg.DBType {
	case "mysql":
		a.repository = mysql.New(a.logger, mysql.Options{
			DBIP:           cfg.MySQL.DBIP,
			User:           cfg.MySQL.DBUser,
			Pass:           cfg.MySQL.DBPass,
			DBName:         cfg.MySQL.DBName,
			ConnRetryCount: cfg.MySQL.DBRetryCount,
		})
	case "":
		a.logger.Fatal("missing db type")
	default:
		a.logger.Fatalf("db type not recognized: %s", cfg.DBType)
	}
}

// Start the application
func (a *Api) Start(cfg *config.Configuration) error {
	//
	a.logger.Println("running the application")

	//
	a.Ctx, a.Cancel = context.WithCancel(context.Background())

	go func() {
		sd := make(chan os.Signal, 1)
		signal.Notify(sd, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sd
		a.logger.Printf("Turn off signal %s", sig)
		defer a.Stop()
	}()

	// start the aplication
	a.setRepo(cfg)
	// Add the necessary things to the biz layer
	a.bizLayer = biz.New(a.logger, a.repository)
	// add the rabbitmq instance
	a.bizLayer.SetRabbit(rabbitmq.NewRabbit(a.logger, cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host))

	// start the biz
	if err := a.bizLayer.Start(); err != nil {
		a.logger.Fatalf("Error stating biz layer: %#v", err)
	}
	// add the apis
	a.apis = api.New(a.logger, cfg.Port, a.bizLayer)
	// run the apis after start the biz layer
	if err := a.apis.Run(a.Ctx); err != nil {
		return fmt.Errorf("Error running api: %#v", err)
	}
	<-a.Ctx.Done()
	return nil
}

// Stop the apis and the biz layer
func (a *Api) Stop() {
	errapis := a.apis.TurnOff()
	if errapis != nil {
		a.logger.Fatalf("Error turning off apis: %+v", errapis)
	}
	errbiz := a.bizLayer.Stop()
	if errbiz != nil {
		a.logger.Fatalf("Error stoping the biz layer: %+v", errbiz)
	}
	// execute cancel context to exit
	defer a.Cancel()
}
