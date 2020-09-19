package ue_message

import (
	"sync"
)

var UeChannel chan HandlerMessage
var mtx sync.Mutex

const (
	MaxChannel int = 100000
)

func init() {
	// init Pool
	UeChannel = make(chan HandlerMessage, MaxChannel)
}

func SendMessage(msg HandlerMessage) {
	mtx.Lock()
	UeChannel <- msg
	mtx.Unlock()
}
