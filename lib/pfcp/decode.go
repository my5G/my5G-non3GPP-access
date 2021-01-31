package pfcp

import (
	"fmt"

	"free5gc/lib/tlv"
)

func (m *Message) Unmarshal(data []byte) error {
	if err := m.Header.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("pfcp: unmarshal msg failed: %s", err)
	}

	// Check Message Length field in header
	if int(m.Header.MessageLength) != len(data)-4 {
		return fmt.Errorf("Incorrect Message Length: Expected %d, got %d", m.Header.MessageLength, len(data)-4)
	}
	switch m.Header.MessageType {
	case PFCP_HEARTBEAT_REQUEST:
		Body := HeartbeatRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_HEARTBEAT_RESPONSE:
		Body := HeartbeatResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_PFD_MANAGEMENT_REQUEST:
		Body := PFCPPFDManagementRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_PFD_MANAGEMENT_RESPONSE:
		Body := PFCPPFDManagementResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_SETUP_REQUEST:
		Body := PFCPAssociationSetupRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_SETUP_RESPONSE:
		Body := PFCPAssociationSetupResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_UPDATE_REQUEST:
		Body := PFCPAssociationUpdateRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_UPDATE_RESPONSE:
		Body := PFCPAssociationUpdateResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_RELEASE_REQUEST:
		Body := PFCPAssociationReleaseRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_RELEASE_RESPONSE:
		Body := PFCPAssociationReleaseResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_NODE_REPORT_REQUEST:
		Body := PFCPNodeReportRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_NODE_REPORT_RESPONSE:
		Body := PFCPNodeReportResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_SET_DELETION_REQUEST:
		Body := PFCPSessionSetDeletionRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_SET_DELETION_RESPONSE:
		Body := PFCPSessionSetDeletionResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_ESTABLISHMENT_REQUEST:
		Body := PFCPSessionEstablishmentRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_ESTABLISHMENT_RESPONSE:
		Body := PFCPSessionEstablishmentResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_MODIFICATION_REQUEST:
		Body := PFCPSessionModificationRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_MODIFICATION_RESPONSE:
		Body := PFCPSessionModificationResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_DELETION_REQUEST:
		Body := PFCPSessionDeletionRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_DELETION_RESPONSE:
		Body := PFCPSessionDeletionResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_REPORT_REQUEST:
		Body := PFCPSessionReportRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_REPORT_RESPONSE:
		Body := PFCPSessionReportResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	default:
		return fmt.Errorf("pfcp: unmarshal msg type %d not supported", m.Header.MessageType)
	}
	return nil
}
