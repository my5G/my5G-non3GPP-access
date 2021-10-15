package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalNodeReportType(t *testing.T) {
	testData := NodeReportType{
		Upfr: true,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{1}, buf)
}

func TestUnmarshalNodeReportType(t *testing.T) {
	buf := []byte{1}
	var testData NodeReportType
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := NodeReportType{
		Upfr: true,
	}
	assert.Equal(t, expectData, testData)
}
