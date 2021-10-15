package pfcpType

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalRedirectInformation(t *testing.T) {
	testData := RedirectInformation{
		RedirectAddressType:         2,
		RedirectServerAddressLength: 13,
		RedirectServerAddress:       []byte("free5gc.local"),
	}
	buf, err := testData.MarshalBinary()

	assert.Nil(t, err)
	assert.Equal(t, []byte{2, 0, 13, 102, 114, 101, 101, 53, 103, 99, 46, 108, 111, 99, 97, 108}, buf)
}

func TestUnmarshalRedirectInformation(t *testing.T) {
	buf := []byte{2, 0, 13, 102, 114, 101, 101, 53, 103, 99, 46, 108, 111, 99, 97, 108}
	var testData RedirectInformation
	err := testData.UnmarshalBinary(buf)

	assert.Nil(t, err)
	expectData := RedirectInformation{
		RedirectAddressType:         2,
		RedirectServerAddressLength: 13,
		RedirectServerAddress:       []byte("free5gc.local"),
	}
	assert.Equal(t, expectData, testData)
}
