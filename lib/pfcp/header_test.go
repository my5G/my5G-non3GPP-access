package pfcp_test

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"free5gc/lib/pfcp"
	"testing"
)

const (
	NodeRelatedHeaderHex    = "2005000000000100"
	SessionRelatedHeaderHex = "233400000000000000000001000000F0"
)

var NodeRelatedHeader = pfcp.Header{
	Version:        1,
	S:              0,
	MessageType:    pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
	MessageLength:  0,
	SequenceNumber: 1,
}

var SessionRelatedHeader = pfcp.Header{
	Version:         1,
	MP:              1,
	S:               1,
	MessageType:     pfcp.PFCP_SESSION_MODIFICATION_REQUEST,
	MessageLength:   0,
	SEID:            1,
	SequenceNumber:  0,
	MessagePriority: 15,
}

func TestPFCPHeader_MarshalBinary(t *testing.T) {
	var payload []byte
	var expect []byte

	// Test PFCP Header for Node Related Messages
	payload, _ = NodeRelatedHeader.MarshalBinary()
	expect, _ = hex.DecodeString(NodeRelatedHeaderHex)
	assert.Equal(t, expect, payload)

	// Test PFCP Header for Session Related Messages
	payload, _ = SessionRelatedHeader.MarshalBinary()
	expect, _ = hex.DecodeString(SessionRelatedHeaderHex)
	assert.Equal(t, expect, payload)
}

func TestPFCPHeader_UnmarshalBinary(t *testing.T) {
	var tmpHeader pfcp.Header
	var payload []byte

	// Test PFCP Header for Node Related Messages
	payload, _ = hex.DecodeString(NodeRelatedHeaderHex)
	_ = tmpHeader.UnmarshalBinary(payload)
	assert.Equal(t, NodeRelatedHeader, tmpHeader)

	// Test PFCP Header for Session Related Messages
	payload, _ = hex.DecodeString(SessionRelatedHeaderHex)
	_ = tmpHeader.UnmarshalBinary(payload)
	assert.Equal(t, SessionRelatedHeader, tmpHeader)
}
