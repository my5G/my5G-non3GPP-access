package nasType_test

import (
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeIdentityResponseMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeIdentityResponseMessageIdentityTable = []nasTypeIdentityResponseMessageIdentityData{
	{nas.MsgTypeIdentityResponse, nas.MsgTypeIdentityResponse},
}

func TestNasTypeNewIdentityResponseMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityResponseMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeIdentityResponseMessageIdentity(t *testing.T) {
	a := nasType.NewIdentityResponseMessageIdentity()
	for _, table := range nasTypeIdentityResponseMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
