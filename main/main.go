package main

import (
	"github.com/IvoryRaptor/dragonfly/test"
	"log"
)

func main() {
	k := test.TestKernel{}
	err := k.New()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = k.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	k.WaitStop()
}
