package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalPrecedence(t *testing.T) {
	testData := Precedence{
		PrecedenceValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalPrecedence(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData Precedence
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := Precedence{
		PrecedenceValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
