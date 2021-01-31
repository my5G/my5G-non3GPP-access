package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalPacketDetectionRuleID(t *testing.T) {
	testData := PacketDetectionRuleID{
		RuleId: 256,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{1, 0}, buf)
}

func TestUnmarshalPacketDetectionRuleID(t *testing.T) {
	buf := []byte{1, 0}
	var testData PacketDetectionRuleID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := PacketDetectionRuleID{
		RuleId: 256,
	}
	assert.Equal(t, expectData, testData)
}
