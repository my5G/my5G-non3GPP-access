package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalTimeQuota(t *testing.T) {
	testData := TimeQuota{
		TimeQuotaValue: 12345,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{0, 0, 48, 57}, buf)
}

func TestUnmarshalTimeQuota(t *testing.T) {
	buf := []byte{0, 0, 48, 57}
	var testData TimeQuota
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := TimeQuota{
		TimeQuotaValue: 12345,
	}
	assert.Equal(t, expectData, testData)
}
