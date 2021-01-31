package fsm_test

import (
	"fmt"
	"free5gc/lib/fsm"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Opened fsm.StateType = "Opened"
	Closed fsm.StateType = "Closed"
)

const (
	Open  fsm.EventType = "Open"
	Close fsm.EventType = "Close"
)

func TestState(t *testing.T) {
	s := fsm.NewState(Closed)

	assert.Equal(t, Closed, s.Current(), "Current() failed")
	assert.True(t, s.Is(Closed), "Is() failed")

	s.Set(Opened)

	assert.Equal(t, Opened, s.Current(), "Current() failed")
	assert.True(t, s.Is(Opened), "Is() failed")
}

func TestFSM(t *testing.T) {
	f, err := fsm.NewFSM(fsm.Transitions{
		{Event: Open, From: Closed, To: Opened},
		{Event: Close, From: Opened, To: Closed},
		{Event: Open, From: Opened, To: Opened},
		{Event: Close, From: Closed, To: Closed},
	}, fsm.Callbacks{
		Opened: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
		Closed: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
	})

	s := fsm.NewState(Closed)

	assert.Nil(t, err, "NewFSM() failed")

	assert.Nil(t, f.SendEvent(s, Open, fsm.ArgsType{"TestArg": "test arg"}), "SendEvent() failed")
	assert.Nil(t, f.SendEvent(s, Close, fsm.ArgsType{"TestArg": "test arg"}), "SendEvent() failed")
	assert.True(t, s.Is(Closed), "Transition failed")

	fakeEvent := fsm.EventType("fake event")
	assert.EqualError(t, f.SendEvent(s, fakeEvent, nil),
		fmt.Sprintf("Unknown transition[From: %s, Event: %s]", s.Current(), fakeEvent))
}

func TestFSMInitFail(t *testing.T) {
	duplicateTrans := fsm.Transition{
		Event: Close, From: Opened, To: Closed,
	}
	_, err := fsm.NewFSM(fsm.Transitions{
		{Event: Open, From: Closed, To: Opened},
		duplicateTrans,
		duplicateTrans,
		{Event: Open, From: Opened, To: Opened},
		{Event: Close, From: Closed, To: Closed},
	}, fsm.Callbacks{
		Opened: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
		Closed: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
	})

	assert.EqualError(t, err, fmt.Sprintf("Duplicate transition: %+v", duplicateTrans))

	fakeState := fsm.StateType("fake state")

	_, err = fsm.NewFSM(fsm.Transitions{
		{Event: Open, From: Closed, To: Opened},
		{Event: Close, From: Opened, To: Closed},
		{Event: Open, From: Opened, To: Opened},
		{Event: Close, From: Closed, To: Closed},
	}, fsm.Callbacks{
		Opened: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
		Closed: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
		fakeState: func(state *fsm.State, event fsm.EventType, args fsm.ArgsType) {
			fmt.Printf("event [%+v] at state [%+v]\n", event, state.Current())
		},
	})

	assert.EqualError(t, err, fmt.Sprintf("Unknown state: %+v", fakeState))
}
