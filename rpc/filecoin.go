package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// client implements the `Client` interface.
type LotusClient struct {
	host string
	jwt  string
}

// NewClient returns a new client.
func NewLotusClient(host string, token string) *LotusClient {
	return &LotusClient{
		host: host,
		jwt:  token,
	}
}

func (c *LotusClient) GetMessage(txHash string) (*http.Response, error) {
	requestBody := NewRequestBody("Filecoin.GetMessage")
	requestBody.Params = append(requestBody.Params, txHash)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body %v", err)
	}
	client := http.Client{}
	req, err := http.NewRequest("POST", c.host, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "Bearer "+c.jwt)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct lotus node rpc request %v", err)
	}
	return client.Do(req)
}

// HandleRequest implements the `Client` interface.
func (c *LotusClient) HandleRequest(r *http.Request, data []byte) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", c.host, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "Bearer "+c.jwt)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct lotus node rpc request %v", err)
	}
	return client.Do(req)
}
