package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDLBufferingDuration(t *testing.T) {
	testData := DLBufferingDuration{
		TimerUnit:  2,
		TimerValue: 21,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{85}, buf)
}

func TestUnmarshalDLBufferingDuration(t *testing.T) {
	buf := []byte{85}
	var testData DLBufferingDuration
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DLBufferingDuration{
		TimerUnit:  2,
		TimerValue: 21,
	}
	assert.Equal(t, expectData, testData)
}
