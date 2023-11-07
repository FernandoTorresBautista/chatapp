package command

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"chatapp/app/biz"
	"chatapp/app/client/rabbitmq"
	"chatapp/app/config"
	"chatapp/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Command struct {
	cfg      *config.Configuration
	logger   *log.Logger
	Ctx      context.Context
	Cancel   context.CancelFunc
	bizLayer *biz.Biz
	rabbit   *rabbitmq.Rabbit
	port     string
	server   *http.Server
}

var apiURL string = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"

// NewCommand ...
func NewCommand(logger *log.Logger, cfg *config.Configuration, biz *biz.Biz, rmq *rabbitmq.Rabbit) *Command {
	port := strconv.Itoa(int(cfg.Port))
	return &Command{
		logger:   logger,
		cfg:      cfg,
		bizLayer: biz,
		rabbit:   rmq,
		port:     port,
	}
}

// Start ...
func (cmd *Command) Start() {
	err := cmd.bizLayer.Start()
	if err != nil {
		cmd.logger.Fatalf("can't start biz: %+v", err)
	}
	cmd.logger.Printf("command started modules")
}

// Run ...
func (cmd *Command) Run() {
	cmd.Ctx, cmd.Cancel = context.WithCancel(context.Background())
	go func() {
		sd := make(chan os.Signal, 1)
		signal.Notify(sd, syscall.SIGTERM, syscall.SIGINT)

		sig := <-sd
		cmd.logger.Printf("Turn off signal %s", sig)
		defer cmd.Stop()
	}()

	cmd.logger.Printf("command started modules")
	root := gin.Default()
	root.POST("/", cmd.handleCommand)

	cmd.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", cmd.port),
		Handler: root,
	}

	// use goroutine so that we can leverage the graceful shutdown code.
	go func() {
		cmd.logger.Printf("command started api")
		if err := cmd.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cmd.logger.Fatalf("listen: %s\n", err)
		}
	}()
	<-cmd.Ctx.Done()
}

// Stop ...
func (cmd *Command) Stop() {
	err := cmd.bizLayer.Stop()
	if err != nil {
		cmd.logger.Printf("can't stop biz: %+v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := cmd.server.Shutdown(ctx); err != nil {
		cmd.logger.Printf("server forced to shutdown: %s", err)
	}
	defer cmd.Cancel()
}

func (cmd *Command) handleCommand(c *gin.Context) {
	// Leer el comando desde el cuerpo de la solicitud
	cmm := c.Request.FormValue("command")
	room := c.Request.FormValue("room")
	if room == "" {
		room = "common"
	}
	str := cmd.BootToGetResponseFromTheAPI(cmm, room)
	c.JSON(http.StatusOK, gin.H{"resp": str})
}

// BootToGetResponseFromTheAPI ...
func (cmd *Command) BootToGetResponseFromTheAPI(message, roomName string) string {
	if !strings.HasPrefix(message, "/stock=") {
		cmd.logger.Printf("The message doesn't have the format needed: %s", message)
	}
	pos := strings.Index(message, "/stock=")
	stockCode := message[pos+7:]
	var p []byte
	if stockCode == "" {
		p = []byte(fmt.Sprintln("bot: please add a stock_code"))
		return string(p)
	}
	url := fmt.Sprintf(apiURL, stockCode)
	records, err := utils.GetAPIResponse(url)
	if err != nil {
		cmd.logger.Printf("error getting response from the API: %+v\n", err)
		p = []byte(err.Error())
	} else {
		cmd.logger.Printf("close price (most representative of the day*): %s\n", records[1][6])
		fmt.Printf("close price (most representative of the day*): %s\n", records[1][6])
		if strings.Contains(records[1][6], "N/D") {
			p = []byte(fmt.Sprintf("bot: called api with stock=%s, didn't work, result: %s", stockCode, records[1][6]))
		} else {
			p = []byte(fmt.Sprintf("bot: %s quote is $%s per share", stockCode, records[1][6]))
		}
	}
	// to the queue, this could send the response to every client in the chat
	err = cmd.bizLayer.SendMessage(websocket.TextMessage, roomName, string(p))
	if err != nil {
		cmd.logger.Printf("error sending the message to the queue: %+v", err)
		return fmt.Sprintf("error sending the message to the queue: %+v", err)
	}
	return string(p)
}
