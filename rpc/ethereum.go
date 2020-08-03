package rpc

import (
	"bytes"
	"fmt"
	"net/http"
)

type infuraClient struct {
	url     string
	apiKey  string
	network string
}

// NewInfuraClient returns a new infuraClient.
func NewInfuraClient(apiKey string, network string) Client {
	return &infuraClient{
		url:     fmt.Sprintf("https://%s.infura.io/v3", network),
		apiKey:  apiKey,
		network: network,
	}
}

// HandleRequest implements the `Client` interface.
func (infura *infuraClient) HandleRequest(r *http.Request, data []byte) (*http.Response, error) {
	apiKey := infura.apiKey
	if apiKey == "" {
		return nil, fmt.Errorf("Can't find any infura API keys")
	}
	client := http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", infura.url, apiKey), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to construct Infura post request: %v", err)
	}
	return client.Do(req)
}
