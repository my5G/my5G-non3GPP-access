package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMarshalRemoteGTPUPeer(t *testing.T) {
	subtests := []struct {
		name       string
		testData   RemoteGTPUPeer
		expectData []byte
	}{
		{
			name: "IPv4",
			testData: RemoteGTPUPeer{
				V4:          true,
				Ipv4Address: net.ParseIP("12.34.56.78").To4(),
			},
			expectData: []byte{2, 12, 34, 56, 78},
		},
		{
			name: "IPv6",
			testData: RemoteGTPUPeer{
				V6:          true,
				Ipv6Address: net.ParseIP("2001:db8::68").To16(),
			},
			expectData: []byte{1, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			buf, err := subtest.testData.MarshalBinary()

			assert.Nil(t, err)
			assert.Equal(t, subtest.expectData, buf)
		})
	}
}

func TestUnmarshalRemoteGTPUPeer(t *testing.T) {
	subtests := []struct {
		name       string
		buf        []byte
		expectData RemoteGTPUPeer
	}{
		{
			name: "IPv4",
			buf:  []byte{2, 12, 34, 56, 78},
			expectData: RemoteGTPUPeer{
				V4:          true,
				Ipv4Address: net.ParseIP("12.34.56.78").To4(),
			},
		},
		{
			name: "IPv6",
			buf:  []byte{1, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104},
			expectData: RemoteGTPUPeer{
				V6:          true,
				Ipv6Address: net.ParseIP("2001:db8::68").To16(),
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var testData RemoteGTPUPeer
			err := testData.UnmarshalBinary(subtest.buf)

			assert.Nil(t, err)
			assert.Equal(t, subtest.expectData, testData)
		})
	}
}
