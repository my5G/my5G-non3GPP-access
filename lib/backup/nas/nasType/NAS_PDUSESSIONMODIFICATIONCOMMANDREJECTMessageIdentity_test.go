package nasType_test

import (
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity(t *testing.T) {
	a := nasType.NewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity()
	assert.NotNil(t, a)
}

type nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageType struct {
	in  uint8
	out uint8
}

var nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageTypeTable = []nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageType{
	{nas.MsgTypePDUSessionModificationCommandReject, nas.MsgTypePDUSessionModificationCommandReject},
}

func TestNasTypeGetSetPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageType(t *testing.T) {
	a := nasType.NewPDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentity()
	for _, table := range nasTypePDUSESSIONMODIFICATIONCOMMANDREJECTMessageIdentityMessageTypeTable {
		a.SetMessageType(table.in)
		assert.Equal(t, table.out, a.GetMessageType())
	}
}
