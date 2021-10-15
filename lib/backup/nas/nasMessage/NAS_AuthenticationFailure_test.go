package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"reflect"

	"github.com/stretchr/testify/assert"
)

type nasMessageAuthenticationFailureData struct {
	inExtendedProtocolDiscriminator         uint8
	inSecurityHeader                        uint8
	inSpareHalfOctet                        uint8
	inAuthenticationFailureMessageIdentity  uint8
	in5GMMCause                             nasType.Cause5GMM
	inAuthenticationFailureParameter        nasType.AuthenticationFailureParameter
	outExtendedProtocolDiscriminator        uint8
	outSecurityHeader                       uint8
	outSpareHalfOctet                       uint8
	outAuthenticationFailureMessageIdentity uint8
	out5GMMCause                            nasType.Cause5GMM
	outAuthenticationFailureParameter       nasType.AuthenticationFailureParameter
}

var nasMessageAuthenticationFailureTable = []nasMessageAuthenticationFailureData{
	{
		inExtendedProtocolDiscriminator:        0x01,
		inSecurityHeader:                       0x08,
		inSpareHalfOctet:                       0x01,
		inAuthenticationFailureMessageIdentity: 0x01,
		in5GMMCause:                            nasType.Cause5GMM{0, 0xff},
		inAuthenticationFailureParameter:       nasType.AuthenticationFailureParameter{nasMessage.AuthenticationFailureAuthenticationFailureParameterType, 14, [14]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	},
	{
		inExtendedProtocolDiscriminator:        0x01,
		inSecurityHeader:                       0x08,
		inSpareHalfOctet:                       0x01,
		inAuthenticationFailureMessageIdentity: 0x01,
		in5GMMCause:                            nasType.Cause5GMM{0, 0xff},
		inAuthenticationFailureParameter:       nasType.AuthenticationFailureParameter{0x30, 14, [14]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	},
}

func TestNasTypeNewAuthenticationFailure(t *testing.T) {
	a := nasMessage.NewAuthenticationFailure(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewAuthenticationFailureMessage(t *testing.T) {
	logger.NasMsgLog.Infoln("---Test NAS Message: AuthenticationFailureMessage---")
	for i, table := range nasMessageAuthenticationFailureTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewAuthenticationFailure(0)
		b := nasMessage.NewAuthenticationFailure(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.AuthenticationFailureMessageIdentity.SetMessageType(table.inAuthenticationFailureMessageIdentity)
		a.Cause5GMM = table.in5GMMCause
		a.AuthenticationFailureParameter = nasType.NewAuthenticationFailureParameter(nasMessage.AuthenticationFailureAuthenticationFailureParameterType)
		a.AuthenticationFailureParameter = &table.inAuthenticationFailureParameter

		buff := new(bytes.Buffer)
		a.EncodeAuthenticationFailure(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln("data: ", data)
		b.DecodeAuthenticationFailure(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
