package wire

type InvType int

const (
	InvalidInvType = iota
	TxInvType
)

type InvMessage struct {
	Inventory []InvItem
}

type InvItem struct {
	Type InvType
	Hash []byte
}
