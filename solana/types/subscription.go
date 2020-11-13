package types

type Subscription struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	Result       Result `json:"result"`
	Subscription int    `json:"subscription"`
}

type Result struct {
	Context Context `json:"context"`
	Value   Value   `json:"value"`
}

type Context struct {
	Slot int `json:"slot"`
}

type Value struct {
	Pubkey  string  `json:"pubkey"`
	Account Account `json:"account"`
}

type Account struct {
	Data       string `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int    `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  int    `json:"rentEpoch"`
}
