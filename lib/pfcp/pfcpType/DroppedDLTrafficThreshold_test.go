package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDroppedDLTrafficThreshold(t *testing.T) {
	testData := DroppedDLTrafficThreshold{
		Dlby:                        true,
		Dlpa:                        true,
		DownlinkPackets:             123456789123456789,
		NumberOfBytesOfDownlinkData: 987654321987654321,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{3, 1, 182, 155, 75, 172, 208, 95, 21, 13, 180, 218, 95, 126, 244, 18, 177}, buf)
}

func TestUnmarshalDroppedDLTrafficThreshold(t *testing.T) {
	buf := []byte{3, 1, 182, 155, 75, 172, 208, 95, 21, 13, 180, 218, 95, 126, 244, 18, 177}
	var testData DroppedDLTrafficThreshold
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DroppedDLTrafficThreshold{
		Dlby:                        true,
		Dlpa:                        true,
		DownlinkPackets:             123456789123456789,
		NumberOfBytesOfDownlinkData: 987654321987654321,
	}
	assert.Equal(t, expectData, testData)
}
