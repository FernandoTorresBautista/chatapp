package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"chatapp/app/biz/models"
	"chatapp/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	listNames        = []string{}
	apiURL    string = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

// WebsocketHandler ...
func (api *Apiv1) WebsocketHandler(c *gin.Context) {
	token, _ := c.Cookie("ca-token")
	user := &models.User{}
	err := json.Unmarshal([]byte(token), user)
	if err != nil {
		api.logger.Printf("Error upgrading websocket: %+v\n", err)
		return
	}
	api.logger.Println("user connected: ", user)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		api.logger.Printf("Error upgrading websocket: %+v\n", err)
		return
	}
	defer conn.Close()

	// get the name of the room
	roomName := c.Param("id")

	if api.chatRooms[roomName] == nil {
		api.chatRooms[roomName] = &models.ListWS{
			List: []*websocket.Conn{},
		}
	}

	// add the new client
	api.chatRooms[roomName].List = append(api.chatRooms[roomName].List, conn)

	// create the rabbitmq queue
	err = api.bizLayer.CreateQueue(roomName)
	if err != nil {
		api.logger.Printf("error creating the queue for the room: %+v", err)
		return
	}

	go api.CreateConsumer(roomName)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			api.logger.Printf("error read message: %+v, user %+v\n", err, user)
			delete(api.chatRooms, roomName)
			continue
		}
		message := string(p)

		// if the prefix message start with ""
		if strings.HasPrefix(message, "/stock=") {
			// connect with the api and send the response to the corresponding chat...
			go api.BootToGetResponseFromTheAPI(message, roomName, api.chatRooms[roomName].List)
		}

		sendmsg := user.Name + ":" + string(p)
		// send the message to rabbit, the consumer will send the message to each client
		if err := api.bizLayer.SendMessage(messageType, roomName, sendmsg); err != nil {
			api.logger.Fatalf("error sending the message to the queue: %+v", err)
		}
	}
}

// BootToGetResponseFromTheAPI ...
func (api *Apiv1) BootToGetResponseFromTheAPI(message, roomName string, clients []*websocket.Conn) {
	pos := strings.Index(message, "/stock=")
	stockCode := message[pos+7:]
	url := fmt.Sprintf(apiURL, stockCode)
	records, err := utils.GetAPIResponse(url)
	var p []byte
	if err != nil {
		api.logger.Printf("error getting response from the API: %+v\n", err)
		p = []byte(err.Error())
	} else {
		api.logger.Println("close price (most representative of the day*): ", records[1][6])
		p = []byte(fmt.Sprintf("%s quote is $%s per share", stockCode, records[1][6]))
	}
	err = api.bizLayer.SendMessage(websocket.TextMessage, roomName, string(p))
	if err != nil {
		api.logger.Fatalf("error sending the message to the queue: %+v", err)
		return
	}
}

// CreateConsumer ....
func (api *Apiv1) CreateConsumer(room string) {
	ch, err := api.rabbit.Rmq.Channel()
	if err != nil {
		api.logger.Fatalf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()
	msgs, err := ch.Consume(
		room,  // queue
		room,  // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		api.logger.Fatalf("Failed to create the consumer: %v", err)
		return
	}
	for msg := range msgs {
		// transmit the message to all the clients
		for _, client := range api.chatRooms[room].List {
			if err := client.WriteMessage(websocket.TextMessage, []byte(msg.Body)); err != nil {
				api.logger.Fatalf("error writing messages %+v", err)
				return
			}
		}
	}
}

// JoinCommon ...
func (api *Apiv1) JoinCommon(c *gin.Context) {
	c.HTML(200, "chat.html", gin.H{"chatID": "common"})
}

// GetListOfChats ...
func (api *Apiv1) GetListOfChats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": listNames})
}

// CreateNewChat ...
func (api *Apiv1) CreateNewChat(c *gin.Context) {
	chatName := c.Request.PostFormValue("name")
	listNames = append(listNames, chatName)
	c.HTML(200, "chat.html", gin.H{"chatID": chatName})
}

// JoinToChat ....
func (api *Apiv1) JoinToChat(c *gin.Context) {
	chatName := c.Param("name")
	c.HTML(200, "chat.html", gin.H{"chatID": chatName})
}
