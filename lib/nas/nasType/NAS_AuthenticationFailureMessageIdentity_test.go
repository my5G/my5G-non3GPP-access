package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewAuthenticationFailureMessageIdentity(t *testing.T) {
	a := nasType.NewAuthenticationFailureMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypeMessageType struct {
	in  uint8
	out uint8
}

var nasTypeMessageTypeTable = []nasTypeMessageType{
	{0x03, 0x03},
}

func TestNasTypeGetSetMessageType(t *testing.T) {
	a := nasType.NewAuthenticationFailureMessageIdentity()
	for _, table := range nasTypeMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}

type AuthenticationFailureMessageIdentityTestDataTemplate struct {
	in  nasType.AuthenticationFailureMessageIdentity
	out nasType.AuthenticationFailureMessageIdentity
}

var authenticationFailureMessageIdentityTestData = []nasType.AuthenticationFailureMessageIdentity{
	{0x03},
}

var authenticationFailureMessageIdentityExpectedTestData = []nasType.AuthenticationFailureMessageIdentity{
	{0x03},
}

var authenticationFailureMessageIdentityTable = []AuthenticationFailureMessageIdentityTestDataTemplate{
	{authenticationFailureMessageIdentityTestData[0], authenticationFailureMessageIdentityExpectedTestData[0]},
}

func TestNasTypeAuthenticationFailureMessageIdentity(t *testing.T) {

	for i, table := range authenticationFailureMessageIdentityTable {
		t.Logf("Test Cnt:%d", i)
		a := nasType.NewAuthenticationFailureMessageIdentity()

		a.SetMessageType(table.in.GetMessageType())
		assert.Equal(t, table.out.GetMessageType(), a.GetMessageType())
	}
}
