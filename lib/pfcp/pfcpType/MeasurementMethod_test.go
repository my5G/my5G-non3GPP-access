package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalMeasurementMethod(t *testing.T) {
	testData := MeasurementMethod{
		Event: true,
		Volum: false,
		Durat: true,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{5}, buf)
}

func TestUnmarshalMeasurementMethod(t *testing.T) {
	buf := []byte{5}
	var testData MeasurementMethod
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := MeasurementMethod{
		Event: true,
		Volum: false,
		Durat: true,
	}
	assert.Equal(t, expectData, testData)
}
