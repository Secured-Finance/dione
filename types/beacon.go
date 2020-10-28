package types

type BeaconEntry struct {
	Round uint64
	Data  []byte
}

type Randomness []byte

func NewBeaconEntry(round uint64, data []byte) BeaconEntry {
	return BeaconEntry{
		Round: round,
		Data:  data,
	}
}
