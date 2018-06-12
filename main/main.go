package main

import "github.com/IvoryRaptor/dragonfly/test"

func main() {
	t :=test.TestKernel{}
	t.Name = "test"
	t.Config()
	t.Start()
	t.WaitStop()
}
