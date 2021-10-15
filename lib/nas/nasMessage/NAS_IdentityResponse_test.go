package nasMessage_test

import (
	"bytes"

	"github.com/free5gc/nas/logger"

	//"fmt"
	"testing"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"

	"reflect"

	"github.com/stretchr/testify/assert"
)

type nasMessageIdentityResponseData struct {
	inExtendedProtocolDiscriminator   uint8
	inSecurityHeader                  uint8
	inSpareHalfOctet                  uint8
	inIdentityResponseMessageIdentity uint8
	inMobileIdentity                  nasType.MobileIdentity
}

var nasMessageIdentityResponseTable = []nasMessageIdentityResponseData{
	{
		inExtendedProtocolDiscriminator:   0x01,
		inSecurityHeader:                  0x08,
		inSpareHalfOctet:                  0x01,
		inIdentityResponseMessageIdentity: nas.MsgTypeIdentityResponse,
		inMobileIdentity: nasType.MobileIdentity{
			Iei:    0,
			Len:    2,
			Buffer: []uint8{0x01, 0x01},
		},
	},
}

func TestNasTypeNewIdentityResponse(t *testing.T) {
	a := nasMessage.NewIdentityResponse(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewIdentityResponseMessage(t *testing.T) {

	for i, table := range nasMessageIdentityResponseTable {
		logger.NasMsgLog.Infoln("Test Cnt:", i)
		a := nasMessage.NewIdentityResponse(0)
		b := nasMessage.NewIdentityResponse(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.IdentityResponseMessageIdentity.SetMessageType(table.inIdentityResponseMessageIdentity)

		a.MobileIdentity = table.inMobileIdentity

		buff := new(bytes.Buffer)
		a.EncodeIdentityResponse(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		b.DecodeIdentityResponse(&data)
		logger.NasMsgLog.Debugln(data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
