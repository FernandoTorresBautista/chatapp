package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "chatapp/app/api/v1"
	"chatapp/app/biz"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// HTTPIn instance
type HTTPIn struct {
	logger   *log.Logger
	port     uint
	server   *http.Server
	bizLayer *biz.Biz
	botURL   string
}

// New api instance
func New(logger *log.Logger, port uint, bizIn *biz.Biz, botURL string) *HTTPIn {
	return &HTTPIn{
		logger:   logger,
		port:     port,
		bizLayer: bizIn,
		botURL:   botURL,
	}
}

// ApiHandle to handle the apis
type ApiHandle interface {
	Run(ctx context.Context) error
	TurnOff() error
}

// CreateRouter wrap endpoints for database
func (h *HTTPIn) CreateRouter() (root *gin.Engine) {
	root = gin.Default()

	// adding templates, css and js to gin
	root.LoadHTMLGlob("app/api/cover/html-template/*")
	root.Static("/css", "app/api/cover/css")
	root.Static("/js", "app/api/cover/js")

	v1root := root.Group("")
	v1.AddRoutes(h.logger, v1root, h.bizLayer, h.botURL)

	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return root
}

// Run start the server
func (h *HTTPIn) Run(ctx context.Context) error {
	h.logger.Println("restapi.Run...waiting for context to be canceled")
	root := h.CreateRouter()
	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", h.port),
		Handler: root,
	}
	// use goroutine so that we can leverage the graceful shutdown code.
	go func() {
		if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.logger.Fatalf("listen: %s\n", err)
		}
	}()
	return nil
}

// TurnOff close the server
func (h *HTTPIn) TurnOff() error {
	h.logger.Println("http server: shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.Fatalf("server forced to shutdown: %s", err)
	}
	h.logger.Println("http server: shutdown complete")
	return nil
}
