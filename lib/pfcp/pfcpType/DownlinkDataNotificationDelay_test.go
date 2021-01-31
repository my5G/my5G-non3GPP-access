package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDownlinkDataNotificationDelay(t *testing.T) {
	testData := DownlinkDataNotificationDelay{
		DelayValue: 21,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{21}, buf)
}

func TestUnmarshalDownlinkDataNotificationDelay(t *testing.T) {
	buf := []byte{21}
	var testData DownlinkDataNotificationDelay
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := DownlinkDataNotificationDelay{
		DelayValue: 21,
	}
	assert.Equal(t, expectData, testData)
}
