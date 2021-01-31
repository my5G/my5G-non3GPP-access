package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeSpareHalfOctetAndSecurityHeaderType(t *testing.T) {
	a := nasType.NewSpareHalfOctetAndSecurityHeaderType()
	assert.NotNil(t, a)
}

type nasTypeSecurityHeaderTypeAndSpareHalfOctetData struct {
	inSecurityHeader  uint8
	inSpareHalfOctet  uint8
	outSecurityHeader uint8
	outSpareHalfOctet uint8
}

var nasTypeSecurityHeaderTypeAndSpareHalfOctetTable = []nasTypeSecurityHeaderTypeAndSpareHalfOctetData{
	{0x8, 0x1, 0x8, 0x01},
}

func TestNasTypeGetSetSpareHalfOctetAndSecurityHeaderType(t *testing.T) {
	a := nasType.NewSpareHalfOctetAndSecurityHeaderType()
	for _, table := range nasTypeSecurityHeaderTypeAndSpareHalfOctetTable {
		a.SetSecurityHeaderType(table.inSecurityHeader)
		assert.Equal(t, table.outSecurityHeader, a.GetSecurityHeaderType())
		a.SetSpareHalfOctet(table.inSpareHalfOctet)
		assert.Equal(t, table.outSpareHalfOctet, a.GetSpareHalfOctet())
	}
}
