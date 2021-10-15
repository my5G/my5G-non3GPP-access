package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalGBR(t *testing.T) {
	testData := GBR{
		UlGbr: 123456789123,
		DlGbr: 987654321987,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{28, 190, 153, 26, 131, 229, 244, 200, 247, 67}, buf)
}

func TestUnmarshalGBR(t *testing.T) {
	buf := []byte{28, 190, 153, 26, 131, 229, 244, 200, 247, 67}
	var testData GBR
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := GBR{
		UlGbr: 123456789123,
		DlGbr: 987654321987,
	}
	assert.Equal(t, expectData, testData)
}
