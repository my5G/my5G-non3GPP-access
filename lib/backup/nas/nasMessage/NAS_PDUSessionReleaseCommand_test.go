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

type nasMessagePDUSessionReleaseCommandData struct {
	inExtendedProtocolDiscriminator           uint8
	inPDUSessionID                            uint8
	inPTI                                     uint8
	inPDUSessionReleaseCommandMessageIdentity uint8
	inCause5GSM                               nasType.Cause5GSM
	inBackoffTimerValue                       nasType.BackoffTimerValue
	inEAPMessage                              nasType.EAPMessage
	inExtendedProtocolConfigurationOptions    nasType.ExtendedProtocolConfigurationOptions
}

var nasMessagePDUSessionReleaseCommandTable = []nasMessagePDUSessionReleaseCommandData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSSessionManagementMessage,
		inPDUSessionID:                  0x01,
		inPTI:                           0x01,
		inPDUSessionReleaseCommandMessageIdentity: 0x01,
		inCause5GSM: nasType.Cause5GSM{
			Iei:   0,
			Octet: 0x01,
		},
		inBackoffTimerValue: nasType.BackoffTimerValue{
			Iei:   nasMessage.PDUSessionReleaseCommandBackoffTimerValueType,
			Len:   2,
			Octet: 0x01,
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.PDUSessionReleaseCommandEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inExtendedProtocolConfigurationOptions: nasType.ExtendedProtocolConfigurationOptions{
			Iei:    nasMessage.PDUSessionReleaseCommandExtendedProtocolConfigurationOptionsType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewPDUSessionReleaseCommand(t *testing.T) {
	a := nasMessage.NewPDUSessionReleaseCommand(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewPDUSessionReleaseCommandMessage(t *testing.T) {

	for i, table := range nasMessagePDUSessionReleaseCommandTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewPDUSessionReleaseCommand(0)
		b := nasMessage.NewPDUSessionReleaseCommand(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.PDUSessionID.SetPDUSessionID(table.inPDUSessionID)
		a.PTI.SetPTI(table.inPTI)
		a.PDUSESSIONRELEASECOMMANDMessageIdentity.SetMessageType(table.inPDUSessionReleaseCommandMessageIdentity)

		a.Cause5GSM = table.inCause5GSM

		a.BackoffTimerValue = nasType.NewBackoffTimerValue(nasMessage.PDUSessionReleaseCommandBackoffTimerValueType)
		a.BackoffTimerValue = &table.inBackoffTimerValue

		a.EAPMessage = nasType.NewEAPMessage(nasMessage.PDUSessionReleaseCommandEAPMessageType)
		a.EAPMessage = &table.inEAPMessage

		a.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionReleaseCommandExtendedProtocolConfigurationOptionsType)
		a.ExtendedProtocolConfigurationOptions = &table.inExtendedProtocolConfigurationOptions

		buff := new(bytes.Buffer)
		a.EncodePDUSessionReleaseCommand(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodePDUSessionReleaseCommand(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)
		//logger.NasMsgLog.Debugln(a.Cause5GSM)
		//logger.NasMsgLog.Debugln(b.Cause5GSM)
		//fmt.Println(a.ExtendedProtocolConfigurationOptions)
		//fmt.Println(b.ExtendedProtocolConfigurationOptions)
		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
