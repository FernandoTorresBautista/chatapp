package v1

import (
	"log"

	"chatapp/app/biz"
	"chatapp/app/biz/models"
	"chatapp/app/client/rabbitmq"

	"github.com/gin-gonic/gin"
)

// Apiv1 application structure
type Apiv1 struct {
	logger    *log.Logger
	bizLayer  biz.Handle
	chatRooms map[string]*models.ListWS
	rabbit    *rabbitmq.Rabbit
}

// AddRoutes entrypoint to add the routes to the application
func AddRoutes(logger *log.Logger, rg *gin.RouterGroup, bizLayer *biz.Biz) {
	// create the api
	api := Apiv1{
		logger:    logger,
		bizLayer:  bizLayer,
		chatRooms: make(map[string]*models.ListWS),
		rabbit:    bizLayer.Rabbit,
	}

	// create the common chat
	// rg.POST("/register", api.checkComplete, api.Register)
	rg.GET("/", api.checkComplete, func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	rg.GET("/login", api.checkComplete, func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	rg.GET("/register", api.checkComplete, func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	rg.POST("/register", api.checkComplete, api.Register)
	rg.POST("/login", api.checkComplete, api.LogIn)

	rg.GET("/logout", api.Logout)

	// use a middleware to check the user information
	rg.POST("/joincommon", api.verifyUser, api.JoinCommon)
	rg.GET("/chat/:id", api.verifyUser, api.WebsocketHandler)
	rg.POST("/createchat", api.verifyUser, api.CreateNewChat)
	rg.GET("/chatslist", api.verifyUser, api.GetListOfChats)
	rg.POST("/join/:name", api.verifyUser, api.JoinToChat)
}
