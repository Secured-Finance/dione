package filecoin

import (
	"encoding/json"
	"fmt"

	ftypes "github.com/Secured-Finance/dione/rpc/filecoin/types"
	"github.com/Secured-Finance/dione/rpc/types"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var filecoinURL = "https://api.node.glif.io/"

// client implements the `Client` interface.
type LotusClient struct {
	host       string
	httpClient *fasthttp.Client
}

// NewClient returns a new client.
func NewLotusClient() *LotusClient {
	return &LotusClient{
		host:       filecoinURL,
		httpClient: &fasthttp.Client{},
	}
}

func (c *LotusClient) GetBlock(cid string) ([]byte, error) {
	i := ftypes.NewCidParam(cid)
	return c.HandleRequest("Filecoin.ChainGetBlock", i)
}

func (c *LotusClient) GetTipSetByHeight(chainEpoch int64) ([]byte, error) {
	i := make([]interface{}, 0)
	i = append(i, chainEpoch, nil)
	return c.HandleRequest("Filecoin.ChainGetTipSetByHeight", i)
}

func (c *LotusClient) GetTransaction(cid string) (string, error) {
	i := ftypes.NewCidParam(cid)
	resp, err := c.HandleRequest("Filecoin.ChainGetMessage", i)
	return string(resp), err
}

func (c *LotusClient) GetNodeVersion() ([]byte, error) {
	return c.HandleRequest("Filecoin.Version", nil)
}

func (c *LotusClient) GetChainHead() ([]byte, error) {
	return c.HandleRequest("Filecoin.ChainHead", nil)
}

func (c *LotusClient) VerifyCid(cid string) ([]byte, error) {
	i := ftypes.NewCidParam(cid)
	return c.HandleRequest("Filecoin.ChainHasObj", i)
}

// HandleRequest implements the `Client` interface.
func (c *LotusClient) HandleRequest(method string, params []interface{}) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.host)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	requestBody := types.NewRPCRequestBody(method)
	requestBody.Params = append(requestBody.Params, params...)
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body %v", err)
	}
	req.AppendBody(body)
	resp := fasthttp.AcquireResponse()
	if err = c.httpClient.Do(req, resp); err != nil {
		logrus.Warn("Failed to construct filecoin node rpc request", err)
		return nil, err
	}
	bodyBytes := resp.Body()
	logrus.Info(string(bodyBytes))
	return bodyBytes, nil
}
