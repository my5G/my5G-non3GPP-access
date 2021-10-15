package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalURSEQN(t *testing.T) {
	testData := URSEQN{
		UrseqnValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalURSEQN(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData URSEQN
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := URSEQN{
		UrseqnValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
