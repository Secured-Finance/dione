package models

type Message struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
	From    string                 `json:"-"`
}
