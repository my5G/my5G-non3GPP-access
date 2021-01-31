package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalMeasurementPeriod(t *testing.T) {
	testData := MeasurementPeriod{
		MeasurementPeriod: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalMeasurementPeriod(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData MeasurementPeriod
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := MeasurementPeriod{
		MeasurementPeriod: 12345,
	}
	assert.Equal(t, expectData, testData)
}
