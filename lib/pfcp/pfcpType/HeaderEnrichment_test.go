package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalHeaderEnrichment(t *testing.T) {
	testData := HeaderEnrichment{
		HeaderType:               0,
		LengthOfHeaderFieldName:  12,
		HeaderFieldName:          []byte("Content-Type"),
		LengthOfHeaderFieldValue: 16,
		HeaderFieldValue:         []byte("application/json"),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 12, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 16, 97, 112, 112, 108, 105, 99,
		97, 116, 105, 111, 110, 47, 106, 115, 111, 110}, buf)
}

func TestUnmarshalHeaderEnrichment(t *testing.T) {
	buf := []byte{0, 12, 67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 16, 97, 112, 112, 108, 105, 99,
		97, 116, 105, 111, 110, 47, 106, 115, 111, 110}
	var testData HeaderEnrichment
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := HeaderEnrichment{
		HeaderType:               0,
		LengthOfHeaderFieldName:  12,
		HeaderFieldName:          []byte("Content-Type"),
		LengthOfHeaderFieldValue: 16,
		HeaderFieldValue:         []byte("application/json"),
	}
	assert.Equal(t, expectData, testData)
}
