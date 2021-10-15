package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageSecurityProtected5GSNASMessageData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet                uint8
	inMessageAuthenticationCode     nasType.MessageAuthenticationCode
	inSequenceNumber                nasType.SequenceNumber
	inPlain5GSNASMessage            nasType.Plain5GSNASMessage
}

var nasMessageSecurityProtected5GSNASMessageTable = []nasMessageSecurityProtected5GSNASMessageData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inMessageAuthenticationCode: nasType.MessageAuthenticationCode{
			Octet: [4]uint8{0x01, 0x01, 0x01, 0x01},
		},
		inSequenceNumber: nasType.SequenceNumber{
			Octet: 0x01,
		},
		inPlain5GSNASMessage: nasType.Plain5GSNASMessage{},
	},
}

func TestNasTypeNewSecurityProtected5GSNASMessage(t *testing.T) {
	a := nasMessage.NewSecurityProtected5GSNASMessage(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewSecurityProtected5GSNASMessageMessage(t *testing.T) {

	for i, table := range nasMessageSecurityProtected5GSNASMessageTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewSecurityProtected5GSNASMessage(0)
		b := nasMessage.NewSecurityProtected5GSNASMessage(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)

		a.MessageAuthenticationCode = table.inMessageAuthenticationCode
		a.SequenceNumber = table.inSequenceNumber
		a.Plain5GSNASMessage = table.inPlain5GSNASMessage

		buff := new(bytes.Buffer)
		a.EncodeSecurityProtected5GSNASMessage(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeSecurityProtected5GSNASMessage(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}
	}
}
