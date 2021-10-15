package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalOffendingIE(t *testing.T) {
	testData := OffendingIE{
		TypeOfOffendingIe: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{48, 57}, buf)
}

func TestUnmarshalOffendingIE(t *testing.T) {
	buf := []byte{48, 57}
	var testData OffendingIE
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := OffendingIE{
		TypeOfOffendingIe: 12345,
	}
	assert.Equal(t, expectData, testData)
}
