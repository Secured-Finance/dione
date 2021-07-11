package filecoin

import (
	"bytes"
	"encoding/json"
	"fmt"

	ftypes "github.com/Secured-Finance/dione/rpc/filecoin/types"
	"github.com/Secured-Finance/dione/rpc/types"
	"github.com/filecoin-project/go-address"

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

// func (c *LotusClient) GetTransaction(cid string) (string, error) {
// 	i := ftypes.NewCidParam(cid)
// 	resp, err := c.HandleRequest("Filecoin.ChainGetMessage", i)
// 	return string(resp), err
// }

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

// Gets signed transaction from Filecoin and returns SignedTransaction struct in byte slice
func (c *LotusClient) GetTransaction(cid string) ([]byte, error) {
	i := ftypes.NewCidParam(cid)
	bodyBytes, err := c.HandleRequest("Filecoin.ChainReadObj", i)
	if err != nil {
		return nil, fmt.Errorf("Failed to get object information %v", err)
	}
	var response types.RPCResponseBody
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body %v", err)
	}

	var msg ftypes.SignedMessage
	if err := msg.UnmarshalCBOR(bytes.NewReader(response.Result)); err != nil {
		if err := msg.Message.UnmarshalCBOR(bytes.NewReader(response.Result)); err != nil {
			return nil, err
		}
	}

	switch msg.Message.From.Protocol() | msg.Message.To.Protocol() {
	case address.BLS:
		msg.Type = ftypes.MessageTypeBLS
	case address.SECP256K1:
		msg.Type = ftypes.MessageTypeSecp256k1
	default:
		return nil, fmt.Errorf("Address has unsupported protocol %v", err)
	}

	b := new(bytes.Buffer)
	if err := msg.MarshalCBOR(b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
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
	logrus.Tracef("Filecoin RPC reply: %v", string(bodyBytes))
	return bodyBytes, nil
}
