package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMarshalOuterHeaderCreation(t *testing.T) {
	subtests := []struct {
		name       string
		testData   OuterHeaderCreation
		expectData []byte
	}{
		{
			name: "GTP-U/UDP/IPv4",
			testData: OuterHeaderCreation{
				OuterHeaderCreationDescription: OuterHeaderCreationGtpUUdpIpv4,
				Teid:                           12345,
				Ipv4Address:                    net.ParseIP("12.34.56.78").To4(),
			},
			expectData: []byte{1, 0, 0, 0, 48, 57, 12, 34, 56, 78},
		},
		{
			name: "UDP/IPv6",
			testData: OuterHeaderCreation{
				OuterHeaderCreationDescription: OuterHeaderCreationUdpIpv6,
				Ipv6Address:                    net.ParseIP("2001:db8::68").To16(),
				PortNumber:                     12345,
			},
			expectData: []byte{8, 0, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 48, 57},
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

func TestUnmarshalOuterHeaderCreation(t *testing.T) {
	subtests := []struct {
		name       string
		buf        []byte
		expectData OuterHeaderCreation
	}{
		{
			name: "GTP-U/UDP/IPv4",
			buf:  []byte{1, 0, 0, 0, 48, 57, 12, 34, 56, 78},
			expectData: OuterHeaderCreation{
				OuterHeaderCreationDescription: OuterHeaderCreationGtpUUdpIpv4,
				Teid:                           12345,
				Ipv4Address:                    net.ParseIP("12.34.56.78").To4(),
			},
		},
		{
			name: "UDP/IPv6",
			buf:  []byte{8, 0, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 48, 57},
			expectData: OuterHeaderCreation{
				OuterHeaderCreationDescription: OuterHeaderCreationUdpIpv6,
				Ipv6Address:                    net.ParseIP("2001:db8::68").To16(),
				PortNumber:                     12345,
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var testData OuterHeaderCreation
			err := testData.UnmarshalBinary(subtest.buf)

			assert.Nil(t, err)
			assert.Equal(t, subtest.expectData, testData)
		})
	}
}
