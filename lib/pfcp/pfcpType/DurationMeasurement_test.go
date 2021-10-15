package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDurationMeasurement(t *testing.T) {
	testData := DurationMeasurement{
		DurationValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalDurationMeasurement(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData DurationMeasurement
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DurationMeasurement{
		DurationValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
