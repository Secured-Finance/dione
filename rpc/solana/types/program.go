package types

// SOLANA PROGRAM SCHEME

type Data struct {
	Parsed  Parsed `json:"parsed"`
	Program string `json:"program"`
	Space   int    `json:"space"`
}

type AuthorizedVoters struct {
	AuthorizedVoter string `json:"authorizedVoter"`
	Epoch           int    `json:"epoch"`
}
type EpochCredits struct {
	Credits         string `json:"credits"`
	Epoch           int    `json:"epoch"`
	PreviousCredits string `json:"previousCredits"`
}
type LastTimestamp struct {
	Slot      int `json:"slot"`
	Timestamp int `json:"timestamp"`
}
type Votes struct {
	ConfirmationCount int `json:"confirmationCount"`
	Slot              int `json:"slot"`
}
type Info struct {
	AuthorizedVoters     []AuthorizedVoters `json:"authorizedVoters"`
	AuthorizedWithdrawer string             `json:"authorizedWithdrawer"`
	Commission           int                `json:"commission"`
	EpochCredits         []EpochCredits     `json:"epochCredits"`
	LastTimestamp        LastTimestamp      `json:"lastTimestamp"`
	NodePubkey           string             `json:"nodePubkey"`
	PriorVoters          []interface{}      `json:"priorVoters"`
	RootSlot             int                `json:"rootSlot"`
	Votes                []Votes            `json:"votes"`
}
type Parsed struct {
	Info Info   `json:"info"`
	Type string `json:"type"`
}
