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

type nasMessageServiceAcceptData struct {
	inExtendedProtocolDiscriminator          uint8
	inSecurityHeader                         uint8
	inSpareHalfOctet                         uint8
	inServiceAcceptMessageIdentity           uint8
	inPDUSessionStatus                       nasType.PDUSessionStatus
	inPDUSessionReactivationResult           nasType.PDUSessionReactivationResult
	inPDUSessionReactivationResultErrorCause nasType.PDUSessionReactivationResultErrorCause
	inEAPMessage                             nasType.EAPMessage
}

var nasMessageServiceAcceptTable = []nasMessageServiceAcceptData{
	{
		inExtendedProtocolDiscriminator: nasMessage.Epd5GSMobilityManagementMessage,
		inSecurityHeader:                0x01,
		inSpareHalfOctet:                0x01,
		inServiceAcceptMessageIdentity:  nas.MsgTypeServiceAccept,
		inPDUSessionStatus: nasType.PDUSessionStatus{
			Iei:    nasMessage.ServiceAcceptPDUSessionStatusType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPDUSessionReactivationResult: nasType.PDUSessionReactivationResult{
			Iei:    nasMessage.ServiceAcceptPDUSessionReactivationResultType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inPDUSessionReactivationResultErrorCause: nasType.PDUSessionReactivationResultErrorCause{
			Iei:    nasMessage.ServiceAcceptPDUSessionReactivationResultErrorCauseType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
		inEAPMessage: nasType.EAPMessage{
			Iei:    nasMessage.ServiceAcceptEAPMessageType,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewServiceAccept(t *testing.T) {
	a := nasMessage.NewServiceAccept(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewServiceAcceptMessage(t *testing.T) {

	for i, table := range nasMessageServiceAcceptTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewServiceAccept(0)
		b := nasMessage.NewServiceAccept(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.ServiceAcceptMessageIdentity.SetMessageType(table.inServiceAcceptMessageIdentity)

		a.PDUSessionStatus = nasType.NewPDUSessionStatus(nasMessage.ServiceAcceptPDUSessionStatusType)
		a.PDUSessionStatus = &table.inPDUSessionStatus

		a.PDUSessionReactivationResult = nasType.NewPDUSessionReactivationResult(nasMessage.ServiceAcceptPDUSessionReactivationResultType)
		a.PDUSessionReactivationResult = &table.inPDUSessionReactivationResult

		a.PDUSessionReactivationResultErrorCause = nasType.NewPDUSessionReactivationResultErrorCause(nasMessage.ServiceAcceptPDUSessionReactivationResultErrorCauseType)
		a.PDUSessionReactivationResultErrorCause = &table.inPDUSessionReactivationResultErrorCause

		a.EAPMessage = nasType.NewEAPMessage(nasMessage.ServiceAcceptEAPMessageType)
		a.EAPMessage = &table.inEAPMessage

		buff := new(bytes.Buffer)
		a.EncodeServiceAccept(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeServiceAccept(&data)
		logger.NasMsgLog.Debugln("Dncode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}
	}
}
