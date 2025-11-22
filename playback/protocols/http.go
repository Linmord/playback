package protocols

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPProtocol struct{}

func (h *HTTPProtocol) Connect(server string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 0}
	req, _ := http.NewRequest("GET", server, nil)
	req.Header.Set("Accept", "audio/*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (h *HTTPProtocol) Name() string {
	return "http"
}
