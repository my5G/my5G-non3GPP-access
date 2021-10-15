package timer_test

import (
	"fmt"
	"free5gc/lib/timer"
	"testing"
	"time"
)

var c chan string

func sendChannel(msg string) {
	c <- msg
}

func TestTimer(t *testing.T) {
	c = make(chan string)
	go func() {
		for {
			msg := <-c
			fmt.Println(msg)
		}
	}()

	// Timer Expired
	timer.StartTimer(1, func(msg interface{}) {
		sendChannel(msg.(string))
	}, "sucess")

	time.Sleep(1 * time.Second)

	// Timer Stop Before Expired
	timer := timer.StartTimer(1, func(msg interface{}) {
		sendChannel(msg.(string))
	}, "sucess")
	sucess := timer.Stop()
	if !sucess {
		t.Error("Timer has been stopped")
	}
	time.Sleep(3 * time.Second)

}
