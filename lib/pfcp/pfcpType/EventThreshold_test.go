package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalEventThreshold(t *testing.T) {
	testData := EventThreshold{
		EventThreshold: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalEventThreshold(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData EventThreshold
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := EventThreshold{
		EventThreshold: 12345,
	}
	assert.Equal(t, expectData, testData)
}
