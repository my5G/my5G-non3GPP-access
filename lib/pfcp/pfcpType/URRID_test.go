package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalURRID(t *testing.T) {
	testData := URRID{
		UrrIdValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalURRID(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData URRID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := URRID{
		UrrIdValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
