package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalDLBufferingSuggestedPacketCount(t *testing.T) {
	subtests := []struct {
		name       string
		testData   DLBufferingSuggestedPacketCount
		expectData []byte
	}{
		{
			name: "1 octet length",
			testData: DLBufferingSuggestedPacketCount{
				PacketCountValue: 21,
			},
			expectData: []byte{21},
		},
		{
			name: "2 octet length",
			testData: DLBufferingSuggestedPacketCount{
				PacketCountValue: 277,
			},
			expectData: []byte{1, 21},
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

func TestUnmarshalDLBufferingSuggestedPacketCount(t *testing.T) {
	subtests := []struct {
		name       string
		buf        []byte
		expectData DLBufferingSuggestedPacketCount
	}{
		{
			name: "1 octet length",
			buf:  []byte{21},
			expectData: DLBufferingSuggestedPacketCount{
				PacketCountValue: 21,
			},
		},
		{
			name: "2 octet length",
			buf:  []byte{1, 21},
			expectData: DLBufferingSuggestedPacketCount{
				PacketCountValue: 277,
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			var testData DLBufferingSuggestedPacketCount
			err := testData.UnmarshalBinary(subtest.buf)

			assert.Nil(t, err)
			assert.Equal(t, subtest.expectData, testData)
		})
	}
}
