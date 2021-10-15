package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalActivatePredefinedRules(t *testing.T) {
	testData := ActivatePredefinedRules{
		PredefinedRulesName: []byte("test"),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{116, 101, 115, 116}, buf)
}

func TestUnmarshalActivatePredefinedRules(t *testing.T) {
	buf := []byte{116, 101, 115, 116}
	var testData ActivatePredefinedRules
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := ActivatePredefinedRules{
		PredefinedRulesName: []byte("test"),
	}
	assert.Equal(t, expectData, testData)
}
