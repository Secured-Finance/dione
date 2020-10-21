package main

import (
	"github.com/Secured-Finance/dione/node"
	"github.com/ipfs/go-log"
)

func main() {
	err := node.Start()
	if err != nil {
		log.Logger("node").Panic(err)
	}
}
