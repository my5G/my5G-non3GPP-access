package nasType_test

import (
	"free5gc/lib/nas/nasMessage"
	"free5gc/lib/nas/nasType"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNasTypeNewT3346Value(t *testing.T) {
	a := nasType.NewT3346Value(nasMessage.RegistrationRejectT3346ValueType)
	assert.NotNil(t, a)
}

var nasTypeT3346ValueIeiTable = []NasTypeIeiData{
	{nasMessage.RegistrationRejectT3346ValueType, nasMessage.RegistrationRejectT3346ValueType},
}

func TestNasTypeT3346ValueGetSetIei(t *testing.T) {
	a := nasType.NewT3346Value(nasMessage.RegistrationRejectT3346ValueType)
	for _, table := range nasTypeRegistrationAcceptNetworkFeatureSupport5GSTable {
		a.SetIei(table.in)
		assert.Equal(t, table.out, a.GetIei())
	}
}

var nasTypeT3346ValueLenTable = []NasTypeLenuint8Data{
	{2, 2},
}

func TestNasTypeT3346ValueGetSetLen(t *testing.T) {
	a := nasType.NewT3346Value(nasMessage.RegistrationRejectT3346ValueType)
	for _, table := range nasTypeRegistrationAcceptNetworkFeatureSupport5GSLenTable {
		a.SetLen(table.in)
		assert.Equal(t, table.out, a.GetLen())
	}
}

type NasTypeT3346ValueGPRSTimer2ValueData struct {
	in  uint8
	out uint8
}

var nasTypeT3346ValueGPRSTimer2ValueTable = []NasTypeT3346ValueGPRSTimer2ValueData{
	{0x2, 0x2},
}

func TestNasTypeT3346ValueGetSetGPRSTimer2Value(t *testing.T) {
	a := nasType.NewT3346Value(nasMessage.RegistrationRejectT3346ValueType)
	for _, table := range nasTypeT3346ValueGPRSTimer2ValueTable {
		a.SetGPRSTimer2Value(table.in)
		assert.Equal(t, table.out, a.GetGPRSTimer2Value())
	}
}

type testT3346ValueDataTemplate struct {
	in  nasType.T3346Value
	out nasType.T3346Value
}

var T3346ValueTestData = []nasType.T3346Value{
	{nasMessage.RegistrationRejectT3346ValueType, 1, 0x05},
}
var T3346ValueExpectedData = []nasType.T3346Value{
	{nasMessage.RegistrationRejectT3346ValueType, 1, 0x05},
}

var T3346ValueDataTestTable = []testT3346ValueDataTemplate{
	{T3346ValueTestData[0], T3346ValueExpectedData[0]},
}

func TestNasTypeT3346Value(t *testing.T) {
	for _, table := range T3346ValueDataTestTable {
		a := nasType.NewT3346Value(nasMessage.RegistrationRejectT3346ValueType)
		a.SetIei(table.in.Iei)
		a.SetLen(table.in.Len)
		a.SetGPRSTimer2Value(0x05)
		assert.Equal(t, table.out.Iei, a.Iei)
		assert.Equal(t, table.out.Len, a.Len)
		assert.Equal(t, table.out.Octet, a.Octet)
	}
}
