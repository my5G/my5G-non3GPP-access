package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalVolumeQuota(t *testing.T) {
	testData := VolumeQuota{
		Dlvol:          true,
		Ulvol:          true,
		Tovol:          true,
		TotalVolume:    987654321987654321,
		UplinkVolume:   123456789123456789,
		DownlinkVolume: 864197532864197532,
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{7, 13, 180, 218, 95, 126, 244, 18, 177, 1, 182, 155, 75, 172, 208, 95, 21, 11, 254, 63, 19, 210, 35, 179, 156}, buf)
}

func TestUnmarshalVolumeQuota(t *testing.T) {
	buf := []byte{7, 13, 180, 218, 95, 126, 244, 18, 177, 1, 182, 155, 75, 172, 208, 95, 21, 11, 254, 63, 19, 210, 35, 179, 156}
	var testData VolumeQuota
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := VolumeQuota{
		Dlvol:          true,
		Ulvol:          true,
		Tovol:          true,
		TotalVolume:    987654321987654321,
		UplinkVolume:   123456789123456789,
		DownlinkVolume: 864197532864197532,
	}
	assert.Equal(t, expectData, testData)
}
