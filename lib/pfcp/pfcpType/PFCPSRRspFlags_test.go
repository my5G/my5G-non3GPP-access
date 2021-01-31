package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalPFCPSRRspFlags(t *testing.T) {
	testData := PFCPSRRspFlags{
		Drobu: true,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{1}, buf)
}

func TestUnmarshalPFCPSRRspFlags(t *testing.T) {
	buf := []byte{1}
	var testData PFCPSRRspFlags
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := PFCPSRRspFlags{
		Drobu: true,
	}
	assert.Equal(t, expectData, testData)
}
