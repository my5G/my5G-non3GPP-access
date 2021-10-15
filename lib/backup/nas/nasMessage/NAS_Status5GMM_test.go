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

type nasMessageStatus5GMMData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet                uint8
	inStatus5GMMMessageIdentity     uint8
	inCause5GMM                     nasType.Cause5GMM
}

var nasMessageStatus5GMMTable = []nasMessageStatus5GMMData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inStatus5GMMMessageIdentity:     nas.MsgTypeStatus5GMM,
		inCause5GMM: nasType.Cause5GMM{
			Octet: 0x01,
		},
	},
}

func TestNasTypeNewStatus5GMM(t *testing.T) {
	a := nasMessage.NewStatus5GMM(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewStatus5GMMMessage(t *testing.T) {

	for i, table := range nasMessageStatus5GMMTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewStatus5GMM(0)
		b := nasMessage.NewStatus5GMM(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)

		a.STATUSMessageIdentity5GMM.SetMessageType(table.inStatus5GMMMessageIdentity)

		a.Cause5GMM = table.inCause5GMM

		buff := new(bytes.Buffer)
		a.EncodeStatus5GMM(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeStatus5GMM(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}
	}
}
