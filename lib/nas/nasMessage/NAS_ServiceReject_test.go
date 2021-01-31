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

type nasMessageServiceRejectData struct {
	inExtendedProtocolDiscriminator uint8
	inSecurityHeader                uint8
	inSpareHalfOctet                uint8
	inServiceRejectMessageIdentity  uint8
	inCause5GMM                     nasType.Cause5GMM
	inPDUSessionStatus              nasType.PDUSessionStatus
	inT3346Value                    nasType.T3346Value
	inEAPMessage                    nasType.EAPMessage
}

var nasMessageServiceRejectTable = []nasMessageServiceRejectData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inServiceRejectMessageIdentity:  nas.MsgTypeServiceReject,
		inCause5GMM: nasType.Cause5GMM{
			Octet: 0x01,
		},
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.ServiceRejectPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inT3346Value: nasType.T3346Value{
			Iei:   nasMessage.ServiceRejectT3346ValueType,
			Len:   2,
			Octet: 0x01,
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.ServiceRejectEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewServiceReject(t *testing.T) {
	a := nasMessage.NewServiceReject(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewServiceRejectMessage(t *testing.T) {

	for i, table := range nasMessageServiceRejectTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewServiceReject(0)
		b := nasMessage.NewServiceReject(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.ServiceRejectMessageIdentity.SetMessageType(table.inServiceRejectMessageIdentity)

		a.Cause5GMM = table.inCause5GMM

		a.PDUSessionStatus = nasType.NewPDUSessionStatus(nasMessage.ServiceRejectPDUSessionStatusType)
		a.PDUSessionStatus = &table.inPDUSessionStatus

		a.T3346Value = nasType.NewT3346Value(nasMessage.ServiceRejectT3346ValueType)
		a.T3346Value = &table.inT3346Value

		a.EAPMessage = nasType.NewEAPMessage(nasMessage.ServiceRejectEAPMessageType)
		a.EAPMessage = &table.inEAPMessage

		buff := new(bytes.Buffer)
		a.EncodeServiceReject(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeServiceReject(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}
	}
}
