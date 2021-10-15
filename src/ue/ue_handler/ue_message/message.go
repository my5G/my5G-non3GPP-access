package ue_message

import "net"

type HandlerMessage struct {
	Event       Event
	UDPSendInfo *UDPSendInfoGroup 	// used only when Event == EventN1UDPMessage
	Addr      *net.UDPAddr 			// used only when Event == EventN1UDPMessage
	UEInnerIP   string     			// used when Event == EventN1TunnelCPMessage || Event == EventN1TunnelUPMessage
	TEID        uint32            	// used only when Event == EventGTPMessage
	Value       interface{}
}

type UDPSendInfoGroup struct {
	ChannelID int
	Addr      *net.UDPAddr
}
