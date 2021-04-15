package handler

import (
	ike_messageUE "free5gc/src/ue/ike/message"
	"net"
)

/*
 Make a channel for concurrency
 */

func SendIKEMessageToN3IWF( udpConn *net.UDPConn, srcAddr, dstAddr *net.UDPAddr, message *ike_messageUE.IKEMessage) {

	ikeLog.Trace( "[Send IKE] message to N3IWF -> ", dstAddr.IP  )
	ikeLog.Trace("Encoding ...")

	pkt, err := ike_messageUE.Encode(message)
	if err != nil {
		ikeLog.Errorln(err)
		return
	}

	ikeLog.Trace("[Sending ...")
	n, err := udpConn.WriteToUDP(pkt, dstAddr)
	if err != nil {
		ikeLog.Error(err)
	}
	if n != len(pkt) {
		ikeLog.Errorf("Not all of the data is sent. Total length: %d. Sent %d.", len(pkt), n)
		return
	}

}

func SendTCPNASMessageToN3IWF(tcpConnWithN3IWF *net.TCPConn , pkt []byte){

	ikeLog.Trace( "[Send TCP NAS ] message to N3IWF -> ", tcpConnWithN3IWF.RemoteAddr()  )
	ikeLog.Trace("Encoding ...")

	n, err := tcpConnWithN3IWF.Write(pkt)
	if err != nil {
		ikeLog.Error(err)
	}
	if n != len(pkt) {
		ikeLog.Errorf("Not all NAS TCP of the data is sent. Total length: %d. Sent %d.", len(pkt), n)
		return
	}
}