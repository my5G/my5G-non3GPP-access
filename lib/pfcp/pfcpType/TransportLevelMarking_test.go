package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalTransportLevelMarking(t *testing.T) {
	testData := TransportLevelMarking{
		TosTrafficClass: []byte{43, 21},
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{43, 21}, buf)
}

func TestUnmarshalTransportLevelMarking(t *testing.T) {
	buf := []byte{43, 21}
	var testData TransportLevelMarking
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := TransportLevelMarking{
		TosTrafficClass: []byte{43, 21},
	}
	assert.Equal(t, expectData, testData)
}
