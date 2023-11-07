package command

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"chatapp/app/biz"
	"chatapp/app/client/rabbitmq"
	"chatapp/pkg/utils"

	"github.com/stretchr/testify/assert"
)

// to run this test I used rabbitmq with docker:
// docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3.8-management

func TestRunCommandAPI(t *testing.T) {
	logger := utils.GetLogger("chatapp.log", "chatapp:: ")
	cfg := utils.GetConfig()
	biz := biz.New(logger, nil)
	// rabbittmq
	cfg.RabbitMQ.User = "user"
	cfg.RabbitMQ.Password = "password"
	cfg.RabbitMQ.Host = "localhost"

	rmq := rabbitmq.NewRabbit(logger, cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host)
	biz.SetRabbit(rmq)

	cc := NewCommand(logger, cfg, biz, rmq)
	cc.Start()
	time.Sleep(time.Duration(2) * time.Second)

	// start the application
	go cc.Run()
	time.Sleep(time.Duration(2) * time.Second)

	// Define the URL to which you want to make the POST request
	apiURL := "http://localhost:8081/"
	stockTest := "aapl.us"
	// Create a map to store form values
	formData := url.Values{
		"command": {"/stock=" + stockTest},
		"room":    {"common"},
	}
	// Encode the form values
	formDataEncoded := formData.Encode()
	// Create a request body with the encoded form data
	requestBody := bytes.NewBufferString(formDataEncoded)
	// Make the POST request
	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", requestBody)
	if err != nil {
		fmt.Println("Error making the POST request:", err)
		return
	}
	defer resp.Body.Close()
	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned a non-OK status: %d\n", resp.StatusCode)
		return
	}
	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading the response:", err)
		return
	}
	// Process the API response
	logger.Printf("resp: %s", string(responseBody))
	fmt.Printf("resp: %s", string(responseBody))
	assert.Contains(t, string(responseBody), fmt.Sprintf("%s quote is", stockTest))
	cc.Stop()
}
