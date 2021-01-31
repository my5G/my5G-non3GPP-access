package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMarshalUserPlaneIPResourceInformation(t *testing.T) {
	testData := UserPlaneIPResourceInformation{
		Assosi:          true,
		Assoni:          true,
		Teidri:          4,
		V6:              true,
		V4:              true,
		TeidRange:       21,
		Ipv4Address:     net.ParseIP("12.34.56.78").To4(),
		Ipv6Address:     net.ParseIP("2001:db8::68").To16(),
		NetworkInstance: []byte("free5gc.local"),
		SourceInterface: 12,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{115, 21, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 13, 102, 114, 101,
		101, 53, 103, 99, 46, 108, 111, 99, 97, 108, 12}, buf)
}

func TestUnmarshalUserPlaneIPResourceInformation(t *testing.T) {
	buf := []byte{115, 21, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 13, 102, 114, 101,
		101, 53, 103, 99, 46, 108, 111, 99, 97, 108, 12}
	var testData UserPlaneIPResourceInformation
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := UserPlaneIPResourceInformation{
		Assosi:          true,
		Assoni:          true,
		Teidri:          4,
		V6:              true,
		V4:              true,
		TeidRange:       21,
		Ipv4Address:     net.ParseIP("12.34.56.78").To4(),
		Ipv6Address:     net.ParseIP("2001:db8::68").To16(),
		NetworkInstance: []byte("free5gc.local"),
		SourceInterface: 12,
	}

	assert.Equal(t, expectData, testData)
}
