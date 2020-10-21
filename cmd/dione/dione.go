package main

import (
	"github.com/Secured-Finance/dione/node"
	"github.com/sirupsen/logrus"
)

func main() {
	err := node.Start()
	if err != nil {
		logrus.Panic(err)
	}
}
