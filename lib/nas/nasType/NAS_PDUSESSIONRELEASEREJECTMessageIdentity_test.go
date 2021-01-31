package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONRELEASEREJECTMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONRELEASEREJECTMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONRELEASEREJECTMessageIdentity struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONRELEASEREJECTMessageIdentityTable = []nasTypePDUSESSIONRELEASEREJECTMessageIdentity{
	{0x03, 0x03},
}

func TestNasTypePDUSESSIONRELEASEREJECTMessageIdentityGetSetMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONRELEASEREJECTMessageIdentity()
	for _, table := range nasTypePDUSESSIONRELEASEREJECTMessageIdentityTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}

type PDUSESSIONRELEASEREJECTMessageIdentityTestDataTemplate struct {
	in  nasType.PDUSESSIONRELEASEREJECTMessageIdentity
	out nasType.PDUSESSIONRELEASEREJECTMessageIdentity
}

var PDUSESSIONRELEASEREJECTMessageIdentityTestData = []nasType.PDUSESSIONRELEASEREJECTMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASEREJECTMessageIdentityExpectedTestData = []nasType.PDUSESSIONRELEASEREJECTMessageIdentity{
	{0x03},
}

var PDUSESSIONRELEASEREJECTMessageIdentityTable = []PDUSESSIONRELEASEREJECTMessageIdentityTestDataTemplate{
	{PDUSESSIONRELEASEREJECTMessageIdentityTestData[0], PDUSESSIONRELEASEREJECTMessageIdentityExpectedTestData[0]},
}

func TestNasTypePDUSESSIONRELEASEREJECTMessageIdentity(t *testing.T) {

	for _, table := range PDUSESSIONRELEASEREJECTMessageIdentityTable {

		a := nasType.NewPDUSESSIONRELEASEREJECTMessageIdentity()

		a.SetMessageType(table.in.GetMessageType())
		assert.Equal(t, table.out.GetMessageType(), a.GetMessageType())
	}
}
