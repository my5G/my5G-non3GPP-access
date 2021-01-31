package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalReportingTriggers(t *testing.T) {
	testData := ReportingTriggers{
		Liusa: true,
		Droth: false,
		Stopt: true,
		Start: false,
		Quhti: true,
		Timth: false,
		Volth: true,
		Perio: false,
		Evequ: true,
		Eveth: false,
		Macar: true,
		Envcl: false,
		Timqu: true,
		Volqu: false,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{170, 42}, buf)
}

func TestUnmarshalReportingTriggers(t *testing.T) {
	buf := []byte{170, 42}
	var testData ReportingTriggers
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := ReportingTriggers{
		Liusa: true,
		Droth: false,
		Stopt: true,
		Start: false,
		Quhti: true,
		Timth: false,
		Volth: true,
		Perio: false,
		Evequ: true,
		Eveth: false,
		Macar: true,
		Envcl: false,
		Timqu: true,
		Volqu: false,
	}
	assert.Equal(t, expectData, testData)
}
