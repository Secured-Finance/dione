package rpc

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// client implements the `Client` interface.
type lotusClient struct {
	host string
	jwt  jwt.Token
}

// NewClient returns a new client.
func NewLotusClient(host string, token jwt.Token) Client {
	return &lotusClient{
		host: host,
		jwt:  token,
	}
}

// HandleRequest implements the `Client` interface.
func (c *lotusClient) HandleRequest(r *http.Request, data []byte) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", c.host, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authentication", "Bearer "+c.jwt.Raw)
	if err != nil {
		return nil, fmt.Errorf("Failed to construct lotus node rpc request %v", err)
	}
	return client.Do(req)
}
