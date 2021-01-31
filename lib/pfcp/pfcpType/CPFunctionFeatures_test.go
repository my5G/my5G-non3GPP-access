package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalCPFunctionFeatures(t *testing.T) {
	testData := CPFunctionFeatures{
		SupportedFeatures: CpFunctionFeaturesLoad | CpFunctionFeaturesOvrl,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{3}, buf)
}

func TestUnmarshalCPFunctionFeatures(t *testing.T) {
	buf := []byte{3}
	var testData CPFunctionFeatures
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := CPFunctionFeatures{
		SupportedFeatures: CpFunctionFeaturesLoad | CpFunctionFeaturesOvrl,
	}
	assert.Equal(t, expectData, testData)
}
