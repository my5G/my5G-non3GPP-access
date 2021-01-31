package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMarshalUEIPAddress(t *testing.T) {
	testData := UEIPAddress{
		Ipv6d:                    true,
		Sd:                       true,
		V4:                       true,
		V6:                       true,
		Ipv4Address:              net.ParseIP("12.34.56.78").To4(),
		Ipv6Address:              net.ParseIP("2001:db8::68").To16(),
		Ipv6PrefixDelegationBits: 4,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{15, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 4}, buf)
}

func TestUnmarshalUEIPAddress(t *testing.T) {
	buf := []byte{15, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 4}
	var testData UEIPAddress
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := UEIPAddress{
		Ipv6d:                    true,
		Sd:                       true,
		V4:                       true,
		V6:                       true,
		Ipv4Address:              net.ParseIP("12.34.56.78").To4(),
		Ipv6Address:              net.ParseIP("2001:db8::68").To16(),
		Ipv6PrefixDelegationBits: 4,
	}
	assert.Equal(t, expectData, testData)
}
