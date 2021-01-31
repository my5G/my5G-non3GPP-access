package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDeactivatePredefinedRules(t *testing.T) {
	testData := DeactivatePredefinedRules{
		PredefinedRulesName: []byte("test"),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{116, 101, 115, 116}, buf)
}

func TestUnmarshalDeactivatePredefinedRules(t *testing.T) {
	buf := []byte{116, 101, 115, 116}
	var testData DeactivatePredefinedRules
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DeactivatePredefinedRules{
		PredefinedRulesName: []byte("test"),
	}
	assert.Equal(t, expectData, testData)
}
