package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalLinkedURRID(t *testing.T) {
	testData := LinkedURRID{
		LinkedUrrIdValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalLinkedURRID(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData LinkedURRID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := LinkedURRID{
		LinkedUrrIdValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
