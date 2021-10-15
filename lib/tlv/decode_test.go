package tlv_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/tlv"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	b, _ := hex.DecodeString("000F000F0014000548656C6C6F002800020001FFFF000131FFFF00023232FFFF0003333333FFFF000434343434FFFF00053535353535FFFF0006363636363636FFFF000737373737373737FFFF0008383838383838383800FF0002000100FF0002000200FF0002000300FF0002000400FF00020005")
	testData := TLVTest{}

	err := tlv.Unmarshal(b, &testData)

	assert.Nil(t, err)

	assert.Equal(t, TLVTest{
		StructTest: &StructTest{
			Name:     []byte("Hello"),
			Sequence: 1,
		},
		BinaryMarshalTest: []BinaryMarshalTest{
			{Value: 1},
			{Value: 22},
			{Value: 333},
			{Value: 4444},
			{Value: 55555},
			{Value: 666666},
			{Value: 7777777},
			{Value: 88888888},
		},
		SliceTest: []uint16{
			1, 2, 3, 4, 5,
		},
	}, testData, "encode error")
}
