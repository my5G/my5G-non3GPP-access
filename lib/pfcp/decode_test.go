package pfcp_test

import (
	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpType"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshal(t *testing.T) {

	var testCases = []struct {
		name            string
		expectedMessage pfcp.Message
		inputBuf        []byte
	}{
		{
			name: "Association Setup Request",
			expectedMessage: pfcp.Message{
				Header: pfcp.Header{
					Version:         1,
					MP:              0,
					S:               0,
					MessageType:     pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
					MessageLength:   13,
					SEID:            0,
					SequenceNumber:  1,
					MessagePriority: 0,
				},
				Body: pfcp.PFCPAssociationSetupRequest{
					NodeID: &pfcpType.NodeID{
						NodeIdType:  0,
						NodeIdValue: net.ParseIP("192.188.2.2").To4(),
					},
				},
			},
			inputBuf: []byte{0x20, 0x5, 0x0, 0xd, 0x0, 0x0, 0x1, 0x0, 0x0, 0x3c, 0x0, 0x5, 0x0, 0xc0, 0xbc, 0x2, 0x2},
		},
		{
			name: "Session Establishment Request",
			expectedMessage: pfcp.Message{
				Header: pfcp.Header{
					Version:         1,
					MP:              0,
					S:               1,
					MessageType:     pfcp.PFCP_SESSION_ESTABLISHMENT_REQUEST,
					MessageLength:   44,
					SEID:            0x02,
					SequenceNumber:  1,
					MessagePriority: 0,
				},
				Body: pfcp.PFCPSessionEstablishmentRequest{
					NodeID: &pfcpType.NodeID{
						NodeIdType:  0,
						NodeIdValue: net.ParseIP("192.188.2.2").To4(),
					},
					CreatePDR: []*pfcp.CreatePDR{
						{
							PDRID:              &pfcpType.PacketDetectionRuleID{RuleId: 1},
							Precedence:         &pfcpType.Precedence{PrecedenceValue: 32},
							OuterHeaderRemoval: &pfcpType.OuterHeaderRemoval{OuterHeaderRemovalDescription: pfcpType.OuterHeaderRemovalGtpUUdpIpv4},
						},
					},
				},
			},
			inputBuf: []byte{0x21, 0x32, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x1, 0x0, 0x0, 0x3c, 0x0, 0x5, 0x0, 0xc0, 0xbc, 0x2, 0x2, 0x0, 0x1, 0x0, 0x13, 0x0, 0x38, 0x0, 0x2, 0x0, 0x1, 0x0, 0x1d, 0x0, 0x4, 0x0, 0x0, 0x0, 0x20, 0x0, 0x5f, 0x0, 0x1, 0x0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var msg pfcp.Message
			err := msg.Unmarshal(tc.inputBuf)
			require.Nilf(t, err, "decode %s failed", tc.name)
			require.Equal(t, tc.expectedMessage, msg)
		})
	}
}
