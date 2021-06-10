package types

const TicketRandomnessLookback = 1

// DioneTask represents the values of task computation
type DioneTask struct {
	OriginChain   uint8
	RequestType   string
	RequestParams string
	Payload       []byte
	RequestID     string
}
