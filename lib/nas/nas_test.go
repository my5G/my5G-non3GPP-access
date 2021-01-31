package nas

import (
	"bytes"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var hexString = "7e00560102000021e440b883d63a9f9c56b3703217152eba2010068f241c77748000b2180e54a9760068"

func TestNasGmmMessage(t *testing.T) {
	data, _ := hex.DecodeString(hexString)
	m := NewMessage()
	err := m.GmmMessageDecode(&data)
	assert.Nil(t, err)

	buff := new(bytes.Buffer)
	err = m.GmmMessageEncode(buff)
	assert.Nil(t, err)
}

func TestNasGsmMessage(t *testing.T) {
	data, _ := hex.DecodeString(hexString)
	m := NewMessage()
	err := m.GsmMessageDecode(&data)
	assert.NotNil(t, err)

	buff := new(bytes.Buffer)
	err = m.GsmMessageEncode(buff)
	assert.NotNil(t, err)
}

func TestPlainNas(t *testing.T) {
	data, _ := hex.DecodeString(hexString)
	m := NewMessage()
	err := m.PlainNasDecode(&data)
	assert.Nil(t, err)
	buff, err1 := m.PlainNasEncode()
	assert.Nil(t, err1)
	if !reflect.DeepEqual(data, buff) {
		t.Errorf("Expect : 0x%0x\nOuput: 0x%0x", data, buff)
	}
}
