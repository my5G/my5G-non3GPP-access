package service

import (
	"sync"
	"testing"
	// "free5gc/src/n3iwf/n3iwf"
)

func TestServer(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(2)

	go Run()
	// go n3iwf.Handle()

	wg.Wait()

}
