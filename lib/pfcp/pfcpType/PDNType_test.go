package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalPDNType(t *testing.T) {
	testData := PDNType{
		PdnType: PDNTypeIpv4,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{PDNTypeIpv4}, buf)
}

func TestUnmarshalPDNType(t *testing.T) {
	buf := []byte{PDNTypeIpv4}
	var testData PDNType
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := PDNType{
		PdnType: PDNTypeIpv4,
	}
	assert.Equal(t, expectData, testData)
}
