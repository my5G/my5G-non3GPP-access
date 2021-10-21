package ue_procedures2

import (
	"github.com/free5gc/openapi/models"
)

type milenageTestSet struct {
	K      string
	RAND   string
	SQN    string
	AMF    string
	OP     string
	OPC    string
	F1     string
	F1star string
	F2     string
	F3     string
	F4     string
	F5     string
	F5star string
}

const (
	SUCCESS_CASE                 = "success"
	FAILURE_CASE                 = "failure"
	SUCCESS_SERVING_NETWORK_NAME = "free5gc"
	TESTSET_SERVING_NETWORK_NAME = "WLAN"
)

var TestGenAuthDataTable = make(map[string]*models.AuthenticationInfoRequest)
var MilenageTestSet19 milenageTestSet
var MilenageTestSet18 milenageTestSet

func init() {
	TestGenAuthDataTable[SUCCESS_CASE] = &models.AuthenticationInfoRequest{
		ServingNetworkName: TESTSET_SERVING_NETWORK_NAME,
	}

	// TS 35.208 test set 19
	MilenageTestSet19 = milenageTestSet{
		K:      "5122250214c33e723a5dd523fc145fc0",
		RAND:   "81e92b6c0ee0e12ebceba8d92a99dfa5",
		SQN:    "16f3b3f70fc2",
		AMF:    "c3ab",
		OP:     "c9e8763286b5b9ffbdf56e1297d0887b",
		OPC:    "981d464c7c52eb6e5036234984ad0bcf",
		F1:     "2a5c23d15ee351d5",
		F1star: "62dae3853f3af9d2",
		F2:     "28d7b0f2a2ec3de5",
		F3:     "5349fbe098649f948f5d2e973a81c00f",
		F4:     "9744871ad32bf9bbd1dd5ce54e3e2e5a",
		F5:     "ada15aeb7bb8",
		F5star: "d461bc15475d",
	}

	MilenageTestSet18 = milenageTestSet{
		K:      "b73a90cbcf3afb622dba83c58a8415df",
		RAND:   "b120f1c1a0102a2f507dd543de68281f",
		SQN:    "f1e8a523a36d",
		AMF:    "471b",
		OP:     "b672047e003bb952dca6cb8af0e5b779",
		OPC:    "df0c67868fa25f748b7044c6e7c245b8",
		F1:     "ebd70341bcd415b0",
		F1star: "12359f5d82220c14",
		F2:     "479dd25c20792d63",
		F3:     "66195dbed0313274c5ca7766615fa25e",
		F4:     "66bec707eb2afc476d7408a8f2927b36",
		F5:     "aefdaa5ddd99",
		F5star: "12ec2b87fbb1",
	}

}

