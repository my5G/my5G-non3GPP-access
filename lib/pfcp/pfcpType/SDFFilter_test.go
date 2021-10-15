package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalSDFFilter(t *testing.T) {
	testData := SDFFilter{
		Bid:                     true,
		Fl:                      true,
		Spi:                     true,
		Ttc:                     true,
		Fd:                      true,
		LengthOfFlowDescription: 2,
		FlowDescription:         []byte("Hi"),
		TosTrafficClass:         []byte{43, 21},
		SecurityParameterIndex:  []byte{12, 34, 56, 78},
		FlowLabel:               []byte{7, 65, 43},
		SdfFilterId:             8388608,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{31, 0, 0, 2, 72, 105, 43, 21, 12, 34, 56, 78, 7, 65, 43, 0, 128, 0, 0}, buf)
}

func TestUnmarshalSDFFilter(t *testing.T) {
	buf := []byte{31, 0, 0, 2, 72, 105, 43, 21, 12, 34, 56, 78, 7, 65, 43, 0, 128, 0, 0}
	var testData SDFFilter
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := SDFFilter{
		Bid:                     true,
		Fl:                      true,
		Spi:                     true,
		Ttc:                     true,
		Fd:                      true,
		LengthOfFlowDescription: 2,
		FlowDescription:         []byte("Hi"),
		TosTrafficClass:         []byte{43, 21},
		SecurityParameterIndex:  []byte{12, 34, 56, 78},
		FlowLabel:               []byte{7, 65, 43},
		SdfFilterId:             8388608,
	}
	assert.Equal(t, expectData, testData)
}
