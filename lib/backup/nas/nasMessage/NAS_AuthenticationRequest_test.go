package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
	"reflect"
)

type nasMessageAuthenticationRequestData struct {
	inExtendedProtocolDiscriminator        uint8
	inSecurityHeader                       uint8
	inSpareHalfOctet1                      uint8
	inAuthenticationRequestMessageIdentity uint8
	inTsc                                  uint8
	inNASKeySetIdentifier                  uint8
	inSpareHalfOctet2                      uint8
	inABBA                                 nasType.ABBA
	inAuthenticationParameterRAND          nasType.AuthenticationParameterRAND
	inAuthenticationParameterAUTN          nasType.AuthenticationParameterAUTN
	inEAPMessage                           nasType.EAPMessage
}

var nasMessageAuthenticationRequestTable = []nasMessageAuthenticationRequestData{
	{
		inExtendedProtocolDiscriminator:        0x01,
		inSecurityHeader:                       0x08,
		inSpareHalfOctet1:                      0x01,
		inAuthenticationRequestMessageIdentity: 0x01,
		inTsc:                                  0x01,
		inNASKeySetIdentifier:                  0x07,
		inSpareHalfOctet2:                      0x07,
		inABBA:                                 nasType.ABBA{0, 2, []byte{0x00, 0x00}},
		inAuthenticationParameterRAND:          nasType.AuthenticationParameterRAND{nasMessage.AuthenticationRequestAuthenticationParameterRANDType, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		inAuthenticationParameterAUTN:          nasType.AuthenticationParameterAUTN{nasMessage.AuthenticationRequestAuthenticationParameterAUTNType, 16, [16]uint8{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
		inEAPMessage:                           nasType.EAPMessage{nasMessage.AuthenticationRequestEAPMessageType, 2, []byte{0x00, 0x00}},
	},
}

func TestNasTypeNewAuthenticationRequest(t *testing.T) {
	a := nasMessage.NewAuthenticationRequest(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewAuthenticationRequestMessage(t *testing.T) {
	logger.NasMsgLog.Infoln("---Test NAS Message: AuthenticationRequestMessage---")
	for i, table := range nasMessageAuthenticationRequestTable {
		t.Logf("Test Cnt:%d", i)
		a := nasMessage.NewAuthenticationRequest(0)
		b := nasMessage.NewAuthenticationRequest(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeader)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet1)
		a.AuthenticationRequestMessageIdentity.SetMessageType(table.inAuthenticationRequestMessageIdentity)

		a.ABBA = table.inABBA

		a.AuthenticationParameterRAND = nasType.NewAuthenticationParameterRAND(nasMessage.AuthenticationRequestAuthenticationParameterRANDType)
		a.AuthenticationParameterRAND = &table.inAuthenticationParameterRAND

		a.AuthenticationParameterAUTN = nasType.NewAuthenticationParameterAUTN(nasMessage.AuthenticationRequestAuthenticationParameterAUTNType)
		a.AuthenticationParameterAUTN = &table.inAuthenticationParameterAUTN

		a.EAPMessage = nasType.NewEAPMessage(nasMessage.AuthenticationRequestEAPMessageType)
		a.EAPMessage = &table.inEAPMessage

		buff := new(bytes.Buffer)
		a.EncodeAuthenticationRequest(buff)
		//fmt.Printf("Encode: %x\n", buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		b.DecodeAuthenticationRequest(&data)
		//fmt.Printf("Decode: %x\n", data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}
	}
}
