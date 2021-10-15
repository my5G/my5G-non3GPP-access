package nasMessage_test

import (
	"bytes"
	"free5gc/lib/nas"
	"free5gc/lib/nas/logger"
	"free5gc/lib/nas/nasMessage"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasMessageConfigurationUpdateCompleteData struct {
	inExtendedProtocolDiscriminator              uint8
	inSecurityHeaderType                         uint8
	inSpareHalfOctet                             uint8
	inConfigurationUpdateCompleteMessageIdentity uint8
}

var nasMessageConfigurationUpdateCompleteTable = []nasMessageConfigurationUpdateCompleteData{
	{
		inExtendedProtocolDiscriminator:              nasMessage.Epd5GSSessionManagementMessage,
		inSecurityHeaderType:                         0x01,
		inSpareHalfOctet:                             0x01,
		inConfigurationUpdateCompleteMessageIdentity: nas.MsgTypeConfigurationUpdateComplete,
	},
}

func TestNasTypeNewConfigurationUpdateComplete(t *testing.T) {
	a := nasMessage.NewConfigurationUpdateComplete(0)
	assert.NotNil(t, a)
}

func TestNasTypeNewConfigurationUpdateCompleteMessage(t *testing.T) {

	for i, table := range nasMessageConfigurationUpdateCompleteTable {
		logger.NasMsgLog.Infoln("Test Cnt:", i)
		a := nasMessage.NewConfigurationUpdateComplete(0)
		b := nasMessage.NewConfigurationUpdateComplete(0)
		assert.NotNil(t, a)
		assert.NotNil(t, b)

		a.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(table.inExtendedProtocolDiscriminator)
		a.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(table.inSecurityHeaderType)
		a.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(table.inSpareHalfOctet)
		a.ConfigurationUpdateCompleteMessageIdentity.SetMessageType(table.inConfigurationUpdateCompleteMessageIdentity)

		buff := new(bytes.Buffer)
		a.EncodeConfigurationUpdateComplete(buff)
		logger.NasMsgLog.Debugln("Encode: ", a)

		data := make([]byte, buff.Len())
		buff.Read(data)
		logger.NasMsgLog.Debugln(data)
		b.DecodeConfigurationUpdateComplete(&data)
		logger.NasMsgLog.Debugln("Decode: ", b)

		if reflect.DeepEqual(a, b) != true {
			t.Errorf("Not correct")
		}

	}
}
