package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalReportType(t *testing.T) {
	testData := ReportType{
		Upir: true,
		Erir: false,
		Usar: true,
		Dldr: false,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{10}, buf)
}

func TestUnmarshalReportType(t *testing.T) {
	buf := []byte{10}
	var testData ReportType
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := ReportType{
		Upir: true,
		Erir: false,
		Usar: true,
		Dldr: false,
	}
	assert.Equal(t, expectData, testData)
}
