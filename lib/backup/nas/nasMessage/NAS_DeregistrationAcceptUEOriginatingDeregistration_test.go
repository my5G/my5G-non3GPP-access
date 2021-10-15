package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageDeregistrationAcceptUEOriginatingDeregistrationData struct {
	inExtendedProtocolDiscriminator       uint8
	inSecurityHeaderType                  uint8
	inSpareHalfOctet                      uint8
	inDeregistrationAcceptMessageIdentity uint8
}

var nasMessageDeregistrationAcceptUEOriginatingDeregistrationTable = []nasMessageDeregistrationAcceptUEOriginatingDeregistrationData{
	{
		inExtendedProtocolDiscriminator:       nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                  0x01,
		inSpareHalfOctet:                      0x01,
		inDeregistrationAcceptMessageIdentity: nas.MsgTypeDeregistrationAcceptUEOriginatingDeregistration,
	},
}

func TestNasTypeNewDeregistrationAcceptUEOriginatingDeregistration(t *testing.T) {
	a := nasMessage.NewDeregistrationAcceptUEOriginatingDeregistration(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewDeregistrationAcceptUEOriginatingDeregistrationMessage(t *testing.T) {

	for i, table := range nasMessageDeregistrationAcceptUEOriginatingDeregistrationTable {
		logger.NasMsgLog.Infoln("Test Cnt:", i)
		a := nasMessage.NewDeregistrationAcceptUEOriginatingDeregistration(0)
		b := nasMessage.NewDeregistrationAcceptUEOriginatingDeregistration(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeaderType)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.DeregistrationAcceptMessageIdentity.SetMessageType(table.inDeregistrationAcceptMessageIdentity)

		buff := new(bytes.Buffer)
		a.EncodeDeregistrationAcceptUEOriginatingDeregistration(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeDeregistrationAcceptUEOriginatingDeregistration(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
