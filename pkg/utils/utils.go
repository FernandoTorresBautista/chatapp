package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"chatapp/app/config"
)

// GetAPIResponse make a get request
func GetAPIResponse(url string) ([][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error with the request, body: %s, status code: %d", response.Body.Close().Error(), response.StatusCode)
	}

	reader := csv.NewReader(response.Body)

	// Lee y analiza el contenido del CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// GetLogger return the logger instance
func GetLogger(name, prefix string) *log.Logger {
	logpath := name
	file, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	return log.New(file, prefix, log.LstdFlags|log.Lshortfile)
}

// GetConfig return the current configuration
func GetConfig() *config.Configuration {
	cfg := config.Cfg
	if cfg.Fail {
		fmt.Printf("load configuration failed: %s", cfg.FailMessage)
		os.Exit(1)
	}
	return &cfg
}

// SendBotMessage ...
func SendBotMessage(apiURL, message, room string) error {
	formData := url.Values{
		"command": {message},
		"room":    {room},
	}
	formDataEncoded := formData.Encode()
	requestBody := bytes.NewBufferString(formDataEncoded)
	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", requestBody)
	if err != nil {
		return fmt.Errorf("Error making the POST request: %+v", err)
	}
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bot returned a non-OK status: %d, msg: %s", resp.StatusCode, string(responseBody))
	}
	return nil
}
