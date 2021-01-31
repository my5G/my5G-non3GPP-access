package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDestinationInterface(t *testing.T) {
	testData := DestinationInterface{
		InterfaceValue: DestinationInterfaceCore,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{DestinationInterfaceCore}, buf)
}

func TestUnmarshalDestinationInterface(t *testing.T) {
	buf := []byte{DestinationInterfaceCore}
	var testData DestinationInterface
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DestinationInterface{
		InterfaceValue: DestinationInterfaceCore,
	}
	assert.Equal(t, expectData, testData)
}
