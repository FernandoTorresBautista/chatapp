package utils

import (
	"encoding/csv"
	"fmt"
	"net/http"
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
