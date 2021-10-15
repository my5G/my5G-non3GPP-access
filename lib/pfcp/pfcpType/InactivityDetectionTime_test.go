package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalInactivityDetectionTime(t *testing.T) {
	testData := InactivityDetectionTime{
		InactivityDetectionTime: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalInactivityDetectionTime(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData InactivityDetectionTime
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := InactivityDetectionTime{
		InactivityDetectionTime: 12345,
	}
	assert.Equal(t, expectData, testData)
}
