package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalGateStatus(t *testing.T) {
	testData := GateStatus{
		UlGate: 1,
		DlGate: 1,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{5}, buf)
}

func TestUnmarshalGateStatus(t *testing.T) {
	buf := []byte{5}
	var testData GateStatus
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := GateStatus{
		UlGate: 1,
		DlGate: 1,
	}
	assert.Equal(t, expectData, testData)
}
