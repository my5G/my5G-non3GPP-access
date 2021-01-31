package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalSourceInterface(t *testing.T) {
	testData := SourceInterface{
		InterfaceValue: SourceInterfaceCore,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{SourceInterfaceCore}, buf)
}

func TestUnmarshalSourceInterface(t *testing.T) {
	buf := []byte{SourceInterfaceCore}
	var testData SourceInterface
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := SourceInterface{
		InterfaceValue: SourceInterfaceCore,
	}
	assert.Equal(t, expectData, testData)
}
