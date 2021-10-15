package ue_message

type Event int

const (
	EventRegistrationProcedure Event = iota
	EventDeregistrationProcedure
	EventPDUSessionEstablishment

	EventUEConfigurationUpdate
	EventServiceRequest
	EventPDUSessionModification

	EventN1UDPMessage

	EventN1TunnelCPMessage
	EventN1TunnelUPMessage
	EventGTPMessage
	EventTimerSendRanConfigUpdateMessage
)
