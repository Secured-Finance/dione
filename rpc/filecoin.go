package rpc

import (
	"encoding/json"
	"fmt"

	"github.com/Secured-Finance/dione/lib"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var filecoinURL = "https://filecoin.infura.io/"

// client implements the `Client` interface.
type LotusClient struct {
	host          string
	projectID     string
	projectSecret string
}

// NewClient returns a new client.
func NewLotusClient(pID, secret string) *LotusClient {
	return &LotusClient{
		host:          filecoinURL,
		projectID:     pID,
		projectSecret: secret,
	}
}

func (c *LotusClient) GetMessage(txHash string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.host)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "Basic "+lib.BasicAuth(c.projectID, c.projectSecret))
	requestBody := NewRequestBody("Filecoin.ChainGetMessage")
	requestBody.Params = append(requestBody.Params, txHash)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body %v", err)
	}
	req.AppendBody(body)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	if err = client.Do(req, resp); err != nil {
		logrus.Warn("Failed to construct filecoin node rpc request", err)
		return nil, err
	}
	bodyBytes := resp.Body()
	logrus.Info(string(bodyBytes))
	return resp, nil
}
