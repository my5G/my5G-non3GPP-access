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

type nasMessageDeregistrationRequestUEOriginatingDeregistrationData struct {
	inExtendedProtocolDiscriminator        uint8
	inSecurityHeaderType                   uint8
	inSpareHalfOctet                       uint8
	inDeregistrationRequestMessageIdentity uint8
	inNgksiAndDeregistrationType           nasType.NgksiAndDeregistrationType
	inMobileIdentity5GS                    nasType.MobileIdentity5GS
}

var nasMessageDeregistrationRequestUEOriginatingDeregistrationTable = []nasMessageDeregistrationRequestUEOriginatingDeregistrationData{
	{
		inExtendedProtocolDiscriminator:        nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                   0x01,
		inSpareHalfOctet:                       0x01,
		inDeregistrationRequestMessageIdentity: 0x01,
		inNgksiAndDeregistrationType: nasType.NgksiAndDeregistrationType{
			Octet: 0xFF,
		},
		inMobileIdentity5GS: nasType.MobileIdentity5GS{
			Iei:    0,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewDeregistrationRequestUEOriginatingDeregistration(t *testing.T) {
	a := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewDeregistrationRequestUEOriginatingDeregistrationMessage(t *testing.T) {

	for i, table := range nasMessageDeregistrationRequestUEOriginatingDeregistrationTable {
		logger.NasMsgLog.Infoln("Test Cnt:", i)
		a := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(0)
		b := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeaderType)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.DeregistrationRequestMessageIdentity.SetMessageType(table.inDeregistrationRequestMessageIdentity)

		a.NgksiAndDeregistrationType = table.inNgksiAndDeregistrationType

		a.MobileIdentity5GS = table.inMobileIdentity5GS

		buff := new(bytes.Buffer)
		a.EncodeDeregistrationRequestUEOriginatingDeregistration(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeDeregistrationRequestUEOriginatingDeregistration(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
