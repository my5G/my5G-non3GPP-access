package nasType_test

import (
	"testing"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasType"

	"github.com/stretchr/testify/assert"
)

type nasTypeConfigurationUpdateCommandMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeConfigurationUpdateCommandMessageIdentityTable = []nasTypeConfigurationUpdateCommandMessageIdentityData{
	{nas.MsgTypeConfigurationUpdateCommand, nas.MsgTypeConfigurationUpdateCommand},
}

func TestNasTypeNewConfigurationUpdateCommandMessageIdentity(t *testing.T) {
	a := nasType.NewConfigurationUpdateCommandMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetConfigurationUpdateCommandMessageIdentity(t *testing.T) {
	a := nasType.NewConfigurationUpdateCommandMessageIdentity()
	for _, table := range nasTypeConfigurationUpdateCommandMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
