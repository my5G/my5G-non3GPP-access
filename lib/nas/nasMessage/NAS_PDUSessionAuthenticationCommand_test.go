package nasMessage_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/logger"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"

	"github.com/stretchr/testify/assert"
)

type nasMessagePDUSessionAuthenticationCommandData struct {
	inExtendedProtocolDiscriminator                  uint8
	inPDUSessionID                                   uint8
	inPTI                                            uint8
	inPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity uint8
	inEAPMessage                                     nasType.EAPMessage
	inExtendedProtocolConfigurationOptions           nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionAuthenticationCommandTable = []nasMessagePDUSessionAuthenticationCommandData{
	{
		inExtendedProtocolDiscriminator: nas.MsgTypePDUSessionAuthenticationCommand,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity: 0x01,
		inEAPMessage: nasType.EAPMessage{
			Iei:    0,
			Len:    4,
			Buffer: []uint8{0x01, 0x01, 0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionAuthenticationCommandExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionAuthenticationCommand(t *testing.T) {
	a := nasMessage.NewPDUSessionAuthenticationCommand(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewPDUSessionAuthenticationCommandMessage(t *testing.T) {

	for i, table := range nasMessagePDUSessionAuthenticationCommandTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewPDUSessionAuthenticationCommand(0)
		b := nasMessage.NewPDUSessionAuthenticationCommand(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.PDUSessionID.SetPDUSessionID(table.inPDUSessionID)
		a.PTI.SetPTI(table.inPTI)
		a.PDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity.SetMessageType(table.inPDUSESSIONAUTHENTICATIONCOMMANDMessageIdentity)

		a.EAPMessage = table.inEAPMessage

		a.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionAuthenticationCommandExtendedProtocolConfigurationOptionsType)
		a.ExtendedProtocolConfigurationOptions = &table.inExtendedProtocolConfigurationOptions

		buff := new(bytes.Buffer)
		a.EncodePDUSessionAuthenticationCommand(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodePDUSessionAuthenticationCommand(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
