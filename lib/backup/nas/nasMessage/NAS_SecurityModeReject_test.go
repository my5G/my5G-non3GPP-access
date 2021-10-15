package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageSecurityModeRejectData struct {
	inExtendedProtocolDiscriminator     uint8
	inSecurityHeader                    uint8
	inSpareHalfOctet                    uint8
	inSecurityModeRejectMessageIdentity uint8
	inCause5GMM                         nasType.Cause5GMM
}

var nasMessageSecurityModeRejectTable = []nasMessageSecurityModeRejectData{
	{
		inExtendedProtocolDiscriminator:     nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                    0x01,
		inSpareHalfOctet:                    0x01,
		inSecurityModeRejectMessageIdentity: nas.MsgTypeSecurityModeReject,
		inCause5GMM: nasType.Cause5GMM{
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewSecurityModeReject(t *testing.T) {
	a := nasMessage.NewSecurityModeReject(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewSecurityModeRejectMessage(t *testing.T) {

	for i, table := range nasMessageSecurityModeRejectTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewSecurityModeReject(0)
		b := nasMessage.NewSecurityModeReject(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.SecurityModeRejectMessageIdentity.SetMessageType(table.inSecurityModeRejectMessageIdentity)

		a.Cause5GMM = table.inCause5GMM

		buff := new(bytes.Buffer)
		a.EncodeSecurityModeReject(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeSecurityModeReject(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}
	}
}
