package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalSubsequentTimeThreshold(t *testing.T) {
	testData := SubsequentTimeThreshold{
		SubsequentTimeThreshold: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalSubsequentTimeThreshold(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData SubsequentTimeThreshold
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := SubsequentTimeThreshold{
		SubsequentTimeThreshold: 12345,
	}
	assert.Equal(t, expectData, testData)
}
