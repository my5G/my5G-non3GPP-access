package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalEventID(t *testing.T) {
	testData := EventID{
		EventId: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalEventID(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData EventID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := EventID{
		EventId: 12345,
	}
	assert.Equal(t, expectData, testData)
}
