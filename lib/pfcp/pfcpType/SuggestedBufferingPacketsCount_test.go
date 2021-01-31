package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalSuggestedBufferingPacketsCount(t *testing.T) {
	testData := SuggestedBufferingPacketsCount{
		PacketCountValue: 21,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{21}, buf)
}

func TestUnmarshalSuggestedBufferingPacketsCount(t *testing.T) {
	buf := []byte{21}
	var testData SuggestedBufferingPacketsCount
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := SuggestedBufferingPacketsCount{
		PacketCountValue: 21,
	}
	assert.Equal(t, expectData, testData)
}
