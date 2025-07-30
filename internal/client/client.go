package client

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func New() *http.Client {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	return client
}
func LoadFile(url string, client *http.Client) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Error executing request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}
	return body, nil
}
