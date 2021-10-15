package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMarshalTimeOfFirstPacket(t *testing.T) {
	testData := TimeOfFirstPacket{
		TimeOfFirstPacket: time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{135, 108, 229, 128}, buf)
}

func TestUnmarshalTimeOfFirstPacket(t *testing.T) {
	buf := []byte{135, 108, 229, 128}
	var testData TimeOfFirstPacket
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := TimeOfFirstPacket{
		TimeOfFirstPacket: time.Date(1972, time.January, 1, 0, 0, 0, 0, time.UTC),
	}
	assert.Equal(t, expectData, testData)
}
