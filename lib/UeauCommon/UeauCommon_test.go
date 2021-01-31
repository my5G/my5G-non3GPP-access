package UeauCommon

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

type TestKDF struct {
	NetworkName string
	SQNxorAK    string
	CK          string
	IK          string
	FC          string
	DerivedKey  string
}

const (
	SUCCESS_CASE = "success"
	FAILURE_CASE = "failure"
)

var TestKDFTable = make(map[string]*TestKDF)

func init() {
	// Only the network name is different, which should yet yield different derived results
	TestKDFTable[SUCCESS_CASE] = &TestKDF{
		NetworkName: "WLAN",
		SQNxorAK:    "bb52e91c747a",
		CK:          "5349fbe098649f948f5d2e973a81c00f",
		IK:          "9744871ad32bf9bbd1dd5ce54e3e2e5a",
		FC:          FC_FOR_CK_PRIME_IK_PRIME_DERIVATION,
		DerivedKey:  "0093962d0dd84aa5684b045c9edffa04ccfc230ca74fcc96c0a5d61164f5a76c",
	}

	TestKDFTable[FAILURE_CASE] = &TestKDF{
		NetworkName: "WLANNNNNNNNNNNNNNN",
		SQNxorAK:    "bb52e91c747a",
		CK:          "5349fbe098649f948f5d2e973a81c00f",
		IK:          "9744871ad32bf9bbd1dd5ce54e3e2e5a",
		FC:          FC_FOR_CK_PRIME_IK_PRIME_DERIVATION,
		DerivedKey:  "0093962d0dd84aa5684b045c9edffa04ccfc230ca74fcc96c0a5d61164f5a76c",
	}
}

func TestGetKDFValue(t *testing.T) {
	var NetworkName, SQNxorAK, CK, IK, FC string
	var P0, P1, ckik, val, dk []byte

	// RFC 5448 Test Vector 1
	// success case
	NetworkName = TestKDFTable[SUCCESS_CASE].NetworkName
	P0 = []byte(NetworkName)

	SQNxorAK = TestKDFTable[SUCCESS_CASE].SQNxorAK
	P1, _ = hex.DecodeString(SQNxorAK)

	CK = TestKDFTable[SUCCESS_CASE].CK
	IK = TestKDFTable[SUCCESS_CASE].IK
	ckik, _ = hex.DecodeString(CK + IK)
	// fmt.Printf("ckik = %x\n", ckik)

	FC = TestKDFTable[SUCCESS_CASE].FC
	val = GetKDFValue(ckik, FC, P0, KDFLen(P0), P1, KDFLen(P1))
	// fmt.Printf("val = %x\n", val)

	dk, _ = hex.DecodeString(TestKDFTable[SUCCESS_CASE].DerivedKey)
	if bytes.Equal(val, dk) {
		fmt.Println("TestGetKDFValue success case PASSED.")
	} else {
		fmt.Println("TestGetKDFValue success case FAILED.")
	}

	// failure case
	NetworkName = TestKDFTable[FAILURE_CASE].NetworkName
	P0 = []byte(NetworkName)

	SQNxorAK = TestKDFTable[FAILURE_CASE].SQNxorAK
	P1, _ = hex.DecodeString(SQNxorAK)

	CK = TestKDFTable[FAILURE_CASE].CK
	IK = TestKDFTable[FAILURE_CASE].IK
	ckik, _ = hex.DecodeString(CK + IK)
	// fmt.Printf("ckik = %x\n", ckik)

	FC = TestKDFTable[FAILURE_CASE].FC
	val = GetKDFValue(ckik, FC, P0, KDFLen(P0), P1, KDFLen(P1))
	// fmt.Printf("val = %x\n", val)

	dk, _ = hex.DecodeString(TestKDFTable[FAILURE_CASE].DerivedKey)
	if !bytes.Equal(val, dk) {
		fmt.Println("TestGetKDFValue failure case PASSED.")
	} else {
		fmt.Println("TestGetKDFValue failure case FAILED.")
	}
}
