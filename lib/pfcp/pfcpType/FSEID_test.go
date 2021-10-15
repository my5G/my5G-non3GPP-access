package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMarshalFSEID(t *testing.T) {
	testData := FSEID{
		V4:          true,
		V6:          true,
		Seid:        123456789123456789,
		Ipv4Address: net.ParseIP("12.34.56.78").To4(),
		Ipv6Address: net.ParseIP("2001:db8::68").To16(),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{3, 1, 182, 155, 75, 172, 208, 95, 21, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104}, buf)
}

func TestUnmarshalFSEID(t *testing.T) {
	buf := []byte{3, 1, 182, 155, 75, 172, 208, 95, 21, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104}
	var testData FSEID
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := FSEID{
		V4:          true,
		V6:          true,
		Seid:        123456789123456789,
		Ipv4Address: net.ParseIP("12.34.56.78").To4(),
		Ipv6Address: net.ParseIP("2001:db8::68").To16(),
	}
	assert.Equal(t, expectData, testData)
}
