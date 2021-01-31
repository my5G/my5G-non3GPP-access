package nasType_test

import (
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeIdentityRequestMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeIdentityRequestMessageIdentityTable = []nasTypeIdentityRequestMessageIdentityData{
	{nas.MsgTypeIdentityRequest, nas.MsgTypeIdentityRequest},
}

func TestNasTypeNewIdentityRequestMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityRequestMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeIdentityRequestMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityRequestMessageIdentity()
	for _, table := range nasTypeIdentityRequestMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
