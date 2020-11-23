package types

type FilParams struct {
	Cid interface{} `json:"/"`
}

func NewCidParam(cid interface{}) []interface{} {
	i := make([]interface{}, 0)
	p := &FilParams{
		Cid: cid,
	}
	i = append(i, p)
	return i
}
