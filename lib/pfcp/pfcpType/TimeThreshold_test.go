package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalTimeThreshold(t *testing.T) {
	testData := TimeThreshold{
		TimeThreshold: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalTimeThreshold(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData TimeThreshold
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := TimeThreshold{
		TimeThreshold: 12345,
	}
	assert.Equal(t, expectData, testData)
}
