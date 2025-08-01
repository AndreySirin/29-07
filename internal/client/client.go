package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
			return
		}
	}()

	content := response.Header.Get("Content-Type")
	if strings.HasPrefix(content, "application/pdf") || strings.HasPrefix(content, "image/jpeg") {
		body, errRead := io.ReadAll(response.Body)
		if errRead != nil {
			return nil, fmt.Errorf("error reading response body: %v", errRead)
		}
		return body, nil
	}
	return nil, fmt.Errorf("unexpected content type: %s", content)
}
