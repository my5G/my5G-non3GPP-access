package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalCause(t *testing.T) {
	testData := Cause{
		CauseValue: CauseRequestRejected,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{CauseRequestRejected}, buf)
}

func TestUnmarshalCause(t *testing.T) {
	buf := []byte{CauseRequestRejected}
	var testData Cause
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := Cause{
		CauseValue: CauseRequestRejected,
	}
	assert.Equal(t, expectData, testData)
}
