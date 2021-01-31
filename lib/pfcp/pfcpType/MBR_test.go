package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalMBR(t *testing.T) {
	testData := MBR{
		UlMbr: 123456789123,
		DlMbr: 987654321987,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{28, 190, 153, 26, 131, 229, 244, 200, 247, 67}, buf)
}

func TestUnmarshalMBR(t *testing.T) {
	buf := []byte{28, 190, 153, 26, 131, 229, 244, 200, 247, 67}
	var testData MBR
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := MBR{
		UlMbr: 123456789123,
		DlMbr: 987654321987,
	}
	assert.Equal(t, expectData, testData)
}
