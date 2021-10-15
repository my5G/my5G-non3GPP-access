package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestMarshalFTEID(t *testing.T) {
	subtests := []struct {
		name       string
		testData   FTEID
		expectData []byte
	}{
		{
			name: "Unset CH",
			testData: FTEID{
				Chid:        false,
				Ch:          false,
				V6:          true,
				V4:          true,
				Teid:        12345,
				Ipv4Address: net.ParseIP("12.34.56.78").To4(),
				Ipv6Address: net.ParseIP("2001:db8::68").To16(),
			},
			expectData: []byte{3, 0, 0, 48, 57, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104},
		},
		{
			name: "Set CH",
			testData: FTEID{
				Chid:     true,
				Ch:       true,
				V6:       false,
				V4:       false,
				ChooseId: 7,
			},
			expectData: []byte{12, 7},
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

func TestUnmarshalFTEID(t *testing.T) {
	subtests := []struct {
		name       string
		buf        []byte
		expectData FTEID
	}{
		{
			name: "Unset CH",
			buf:  []byte{3, 0, 0, 48, 57, 12, 34, 56, 78, 32, 1, 13, 184, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104},
			expectData: FTEID{
				Chid:        false,
				Ch:          false,
				V6:          true,
				V4:          true,
				Teid:        12345,
				Ipv4Address: net.ParseIP("12.34.56.78").To4(),
				Ipv6Address: net.ParseIP("2001:db8::68").To16(),
			},
		},
		{
			name: "Set CH",
			buf:  []byte{12, 7},
			expectData: FTEID{
				Chid:     true,
				Ch:       true,
				V6:       false,
				V4:       false,
				ChooseId: 7,
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var testData FTEID
			err := testData.UnmarshalBinary(subtest.buf)

			assert.Nil(t, err)
			assert.Equal(t, subtest.expectData, testData)
		})
	}
}
