package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewOldPDUSessionID(t *testing.T) {
	a := nasType.NewOldPDUSessionID(nasMessage.ULNASTransportOldPDUSessionIDType)
	assert.NotNil(t, a)

}

var nasTypeOldPDUSessionIDULNASTransportOldPDUSessionIDTypeTable = []NasTypeIeiData{
	{nasMessage.ULNASTransportOldPDUSessionIDType, nasMessage.ULNASTransportOldPDUSessionIDType},
}

func TestNasTypeOldPDUSessionIDGetSetIei(t *testing.T) {
	a := nasType.NewOldPDUSessionID(nasMessage.ULNASTransportOldPDUSessionIDType)
	for _, table := range nasTypeOldPDUSessionIDULNASTransportOldPDUSessionIDTypeTable {
		a.SetIei(table.in)
		assert.Equal(t, table.out, a.GetIei())
	}
}

type nasTypeOldPDUSessionIDPduSessionIdentity2Value struct {
	in  uint8
	out uint8
}

var nasTypeOldPDUSessionIDPduSessionIdentity2ValueTable = []nasTypeOldPDUSessionIDPduSessionIdentity2Value{
	{0xff, 0xff},
}

func TestNasTypeOldPDUSessionIDGetSetOldPDUSessionID(t *testing.T) {
	a := nasType.NewOldPDUSessionID(nasMessage.ULNASTransportOldPDUSessionIDType)
	for _, table := range nasTypeOldPDUSessionIDPduSessionIdentity2ValueTable {
		a.SetOldPDUSessionID(table.in)
		assert.Equal(t, table.out, a.GetOldPDUSessionID())
	}
}

type nasTypeOldPDUSessionID struct {
	inIei                       uint8
	inPduSessionIdentity2Value  uint8
	outIei                      uint8
	outPduSessionIdentity2Value uint8
}

var nasTypeOldPDUSessionIDTable = []nasTypeOldPDUSessionID{
	{nasMessage.ULNASTransportOldPDUSessionIDType, 0xff,
		nasMessage.ULNASTransportOldPDUSessionIDType, 0xff},
}

func TestNasTypeOldPDUSessionID(t *testing.T) {
	a := nasType.NewOldPDUSessionID(nasMessage.ULNASTransportOldPDUSessionIDType)
	for _, table := range nasTypeOldPDUSessionIDTable {
		a.SetIei(table.inIei)
		a.SetOldPDUSessionID(table.inPduSessionIdentity2Value)
		assert.Equal(t, table.outIei, a.GetIei())
		assert.Equal(t, table.outPduSessionIdentity2Value, a.GetOldPDUSessionID())
	}
}
