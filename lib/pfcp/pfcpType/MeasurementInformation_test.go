package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalMeasurementInformation(t *testing.T) {
	testData := MeasurementInformation{
		Radi: true,
		Inam: false,
		Mbqe: true,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{5}, buf)
}

func TestUnmarshalMeasurementInformation(t *testing.T) {
	buf := []byte{5}
	var testData MeasurementInformation
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := MeasurementInformation{
		Radi: true,
		Inam: false,
		Mbqe: true,
	}
	assert.Equal(t, expectData, testData)
}
