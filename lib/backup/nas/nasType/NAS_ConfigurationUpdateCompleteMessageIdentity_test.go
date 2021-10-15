package nasType_test

import (
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeConfigurationUpdateCompleteMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeConfigurationUpdateCompleteMessageIdentityTable = []nasTypeConfigurationUpdateCompleteMessageIdentityData{
	{nas.MsgTypeConfigurationUpdateComplete, nas.MsgTypeConfigurationUpdateComplete},
}

func TestNasTypeNewConfigurationUpdateCompleteMessageIdentity(t *testing.T) {
	a := nasType.NewConfigurationUpdateCompleteMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetConfigurationUpdateCompleteMessageIdentity(t *testing.T) {
	a := nasType.NewConfigurationUpdateCompleteMessageIdentity()
	for _, table := range nasTypeConfigurationUpdateCompleteMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
