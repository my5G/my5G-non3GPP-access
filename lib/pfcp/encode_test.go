package pfcp_test

import (
	"free5gc/lib/pfcp"
	"free5gc/lib/pfcp/pfcpType"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {

	var testCases = []struct {
		name        string
		message     pfcp.Message
		expectedBuf []byte
	}{
		{
			name: "Association Setup Request",
			message: pfcp.Message{
				Header: pfcp.Header{
					Version:         1,
					MP:              0,
					S:               0,
					MessageType:     pfcp.PFCP_ASSOCIATION_SETUP_REQUEST,
					MessageLength:   0,
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
			expectedBuf: []byte{0x20, 0x5, 0x0, 0xd, 0x0, 0x0, 0x1, 0x0, 0x0, 0x3c, 0x0, 0x5, 0x0, 0xc0, 0xbc, 0x2, 0x2},
		},
		{
			name: "Session Establishment Request",
			message: pfcp.Message{
				Header: pfcp.Header{
					Version:         1,
					MP:              0,
					S:               1,
					MessageType:     pfcp.PFCP_SESSION_ESTABLISHMENT_REQUEST,
					MessageLength:   0,
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
			expectedBuf: []byte{0x21, 0x32, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x1, 0x0, 0x0, 0x3c, 0x0, 0x5, 0x0, 0xc0, 0xbc, 0x2, 0x2, 0x0, 0x1, 0x0, 0x13, 0x0, 0x38, 0x0, 0x2, 0x0, 0x1, 0x0, 0x1d, 0x0, 0x4, 0x0, 0x0, 0x0, 0x20, 0x0, 0x5f, 0x0, 0x1, 0x0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := tc.message.Marshal()
			require.Nil(t, err, "encode message", tc.name, "failed")
			require.Equal(t, tc.expectedBuf, buf)
		})
	}
}
