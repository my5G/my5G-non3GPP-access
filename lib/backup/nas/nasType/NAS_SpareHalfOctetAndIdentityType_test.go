package nasType_test

import (
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nasTypeIdentityTypeAndSpareHalfOctetData struct {
	in  uint8
	out uint8
}

var nasTypeIdentityTypeAndSpareHalfOctetTable = []nasTypeIdentityTypeAndSpareHalfOctetData{
	{0x07, 0x07},
}

func TestNasTypeNewSpareHalfOctetAndIdentityType(t *testing.T) {
	a := nasType.NewSpareHalfOctetAndIdentityType()
	assert.NotNil(t, a)
}

func TestNasTypeIdentityTypeAndSpareHalfOctet(t *testing.T) {
	a := nasType.NewSpareHalfOctetAndIdentityType()
	for _, table := range nasTypeIdentityTypeAndSpareHalfOctetTable {
		a.SetTypeOfIdentity(table.in)
		assert.Equal(t, table.out, a.GetTypeOfIdentity())
	}
}
