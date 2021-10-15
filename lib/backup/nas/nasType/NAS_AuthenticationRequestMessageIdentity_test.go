package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeRequestMessageIdentityData struct {
	in  uint8
	out uint8
}

var nasTypeRequestMessageIdentityTable = []nasTypeRequestMessageIdentityData{
	{nasMessage.AuthenticationRequestEAPMessageType, nasMessage.AuthenticationRequestEAPMessageType},
}

func TestNasTypeNewAuthenticationRequestMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationRequestMessageIdentity()
	assert.NotNil(t, a)
}

func TestNasTypeGetSetAuthenticationRequestMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationRequestMessageIdentity()
	for _, table := range nasTypeRequestMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
