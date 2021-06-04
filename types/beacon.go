package types

type BeaconEntry struct {
	Round    uint64
	Data     []byte
	Metadata map[string]interface{}
}

type Randomness []byte

func NewBeaconEntry(round uint64, data []byte, metadata map[string]interface{}) BeaconEntry {
	return BeaconEntry{
		Round:    round,
		Data:     data,
		Metadata: metadata,
	}
}
