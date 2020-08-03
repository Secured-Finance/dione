package rpc

import "net/http"

type Client interface {
	HandleRequest(r *http.Request, data []byte) (*http.Response, error)
}
