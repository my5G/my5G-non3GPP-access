package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDownlinkDataServiceInformation(t *testing.T) {
	testData := DownlinkDataServiceInformation{
		Qfii:                        true,
		Ppi:                         true,
		PagingPolicyIndicationValue: 21,
		Qfi:                         12,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{3, 21, 12}, buf)
}

func TestUnmarshalDownlinkDataServiceInformation(t *testing.T) {
	buf := []byte{3, 21, 12}
	var testData DownlinkDataServiceInformation
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DownlinkDataServiceInformation{
		Qfii:                        true,
		Ppi:                         true,
		PagingPolicyIndicationValue: 21,
		Qfi:                         12,
	}
	assert.Equal(t, expectData, testData)
}
