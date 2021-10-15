package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalPFCPSMReqFlags(t *testing.T) {
	testData := PFCPSMReqFlags{
		Qaurr: true,
		Sndem: false,
		Drobu: true,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{5}, buf)
}

func TestUnmarshalPFCPSMReqFlags(t *testing.T) {
	buf := []byte{5}
	var testData PFCPSMReqFlags
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := PFCPSMReqFlags{
		Qaurr: true,
		Sndem: false,
		Drobu: true,
	}
	assert.Equal(t, expectData, testData)
}
