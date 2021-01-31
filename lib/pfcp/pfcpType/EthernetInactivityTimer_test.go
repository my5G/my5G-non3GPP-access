package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalEthernetInactivityTimer(t *testing.T) {
	testData := EthernetInactivityTimer{
		EthernetInactivityTimer: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalEthernetInactivityTimer(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData EthernetInactivityTimer
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := EthernetInactivityTimer{
		EthernetInactivityTimer: 12345,
	}
	assert.Equal(t, expectData, testData)
}
