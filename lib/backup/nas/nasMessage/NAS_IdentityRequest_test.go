package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"

	//"fmt"
	"free5gc/lib/nas"
	"free5gc/lib/nas/nasMessage"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

type nasMessageIdentityRequestData struct {
	inExtendedProtocolDiscriminator  uint8
	inSecurityHeader                 uint8
	inSpareHalfOctet1                uint8
	inIdentityRequestMessageIdentity uint8
	inIdentityType                   uint8
	inSpareHalfOctet2                uint8
}

var nasMessageIdentityRequestTable = []nasMessageIdentityRequestData{
	{
		inExtendedProtocolDiscriminator:  0x01,
		inSecurityHeader:                 0x08,
		inSpareHalfOctet1:                0x01,
		inIdentityRequestMessageIdentity: nas.MsgTypeIdentityRequest,
		inIdentityType:                   0x01,
		inSpareHalfOctet2:                0x01,
	},
}

func TestNasTypeNewIdentityRequest(t *testing.T) {
	a := nasMessage.NewIdentityRequest(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewIdentityRequestMessage(t *testing.T) {

	for i, table := range nasMessageIdentityRequestTable {
		logger.NasMsgLog.Infoln("Test Cnt:", i)
		a := nasMessage.NewIdentityRequest(0)
		b := nasMessage.NewIdentityRequest(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet1)
		a.IdentityRequestMessageIdentity.SetMessageType(table.inIdentityRequestMessageIdentity)
		a.SpareHalfOctetAndIdentityType.SetTypeOfIdentity(table.inIdentityType)

		buff := new(bytes.Buffer)
		a.EncodeIdentityRequest(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		b.DecodeIdentityRequest(&data)
		logger.NasMsgLog.Debugln(data)
		logger.NasMsgLog.Debugln("Dncode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
