package solana

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Secured-Finance/dione/rpc"
	"github.com/Secured-Finance/dione/solana/types"
	ws "github.com/dgrr/fastws"
	"github.com/shengdoushi/base58"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var solanaAlphabet = base58.BitcoinAlphabet

type SolanaClient struct {
	url string
	ws  string
}

type SubParams struct {
	Encoding string `json:"encoding"`
}

func NewSubParam(encoding string) *SubParams {
	return &SubParams{
		Encoding: encoding,
	}
}

// NewSolanaClient creates a new solana client structure.
func NewSolanaClient() *SolanaClient {
	return &SolanaClient{
		url: "http://devnet.solana.com:8899/",
		ws:  "ws://devnet.solana.com:8900/",
	}
}

func (c *SolanaClient) GetTransaction(txHash string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(c.url)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	requestBody := rpc.NewRequestBody("getConfirmedTransaction")
	requestBody.Params = append(requestBody.Params, txHash, "base58")
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request body %v", err)
	}
	req.AppendBody(body)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	if err = client.Do(req, resp); err != nil {
		logrus.Warn("Failed to construct solana node rpc request", err)
		return nil, err
	}
	bodyBytes := resp.Body()
	logrus.Info(string(bodyBytes))
	return resp, nil
}

func (c *SolanaClient) subsctibeOnProgram(programID string) {
	conn, err := ws.Dial(c.ws)
	if err != nil {
		log.Fatalln("Can't establish connection with Solana websocket: ", err)
	}
	defer conn.Close()

	requestBody := rpc.NewRequestBody("programSubscribe")
	requestBody.Params = append(requestBody.Params, programID)
	p := NewSubParam("jsonParsed")
	requestBody.Params = append(requestBody.Params, p)
	body, err := json.Marshal(requestBody)
	logrus.Info(string(body))
	if err != nil {
		logrus.Errorf("Couldn't unmarshal parameters to request body %v", err)
	}

	subscriptionID, err := conn.Write(body)
	if err != nil {
		logrus.Errorf("Couldn't send a websocket request to Solana node %v", err)
	}
	logrus.Info("Websocket established with ProgramID:", programID)
	logrus.Info("Subscription ID to drop websocket connection:", subscriptionID)

	var msg []byte
	var parsedSub *types.Subscription
	for {
		_, msg, err = conn.ReadMessage(msg[:0])
		if err != nil {
			break
		}
		json.Unmarshal(msg, &parsedSub)
		logrus.Info("Subscription: ", parsedSub)
		// 2) Save data from oracle event in redis cache
		// 3) Start mining of Solana oracle event
	}
}
