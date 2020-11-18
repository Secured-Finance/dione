package types

type TxResponse struct {
	Jsonrpc string   `json:"jsonrpc"`
	Result  TxResult `json:"result"`
	ID      int      `json:"id"`
}

type TxStatus struct {
	Ok interface{} `json:"Ok"`
}
type TxMeta struct {
	Err               interface{}   `json:"err"`
	Fee               int           `json:"fee"`
	InnerInstructions []interface{} `json:"innerInstructions"`
	LogMessages       []interface{} `json:"logMessages"`
	PostBalances      []interface{} `json:"postBalances"`
	PreBalances       []interface{} `json:"preBalances"`
	Status            TxStatus      `json:"status"`
}
type TxHeader struct {
	NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	NumRequiredSignatures       int `json:"numRequiredSignatures"`
}
type TxInstructions struct {
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
	ProgramIDIndex int    `json:"programIdIndex"`
}
type Message struct {
	AccountKeys     []string         `json:"accountKeys"`
	Header          TxHeader         `json:"header"`
	Instructions    []TxInstructions `json:"instructions"`
	RecentBlockhash string           `json:"recentBlockhash"`
}
type Transaction struct {
	Message    Message  `json:"message"`
	Signatures []string `json:"signatures"`
}
type TxResult struct {
	Meta        TxMeta      `json:"meta"`
	Slot        int         `json:"slot"`
	Transaction Transaction `json:"transaction"`
}
