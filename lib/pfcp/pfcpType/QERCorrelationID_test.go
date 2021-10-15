package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalQERCorrelationID(t *testing.T) {
	testData := QERCorrelationID{
		QerCorrelationIdValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalQERCorrelationID(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData QERCorrelationID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := QERCorrelationID{
		QerCorrelationIdValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
