package handler

import (
	ike_message "free5gc/src/n3iwf/ike/message"
	"net"
)

func SendIKEMessageToUE(udpConn *net.UDPConn, srcAddr, dstAddr *net.UDPAddr, message *ike_message.IKEMessage) {
	ikeLog.Trace("[IKE] Send IKE message to UE")
	ikeLog.Trace("[IKE] Encoding...")
	pkt, err := ike_message.Encode(message)
	if err != nil {
		ikeLog.Errorln(err)
		return
	}
	// As specified in RFC 7296 section 3.1, the IKE message send from/to UDP port 4500
	// should prepend a 4 bytes zero
	if srcAddr.Port == 4500 {
		prependZero := make([]byte, 4)
		pkt = append(prependZero, pkt...)
	}

	ikeLog.Trace("[IKE] Sending...")
	n, err := udpConn.WriteToUDP(pkt, dstAddr)
	if err != nil {
		ikeLog.Error(err)
		return
	}
	if n != len(pkt) {
		ikeLog.Errorf("Not all of the data is sent. Total length: %d. Sent: %d.", len(pkt), n)
		return
	}
}
