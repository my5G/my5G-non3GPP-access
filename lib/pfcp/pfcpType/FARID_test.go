package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalFARID(t *testing.T) {
	testData := FARID{
		FarIdValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalFARID(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData FARID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := FARID{
		FarIdValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
