package types

//	All ethereum values uses bit.Int
//  to calculate Wei values in GWei:
//	new(big.Int).Mul(value, big.NewInt(params.GWei))
const (
	Wei   = 1
	GWei  = 1e9
	Ether = 1e18
)
