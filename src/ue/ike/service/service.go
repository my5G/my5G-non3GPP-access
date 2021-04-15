package service

import (
	"fmt"
	"free5gc/src/n3iwf/ike"
	"net"
	"errors"

	"free5gc/src/ue/context"
	"free5gc/src/ue/logger"
	"github.com/sirupsen/logrus"
)

/*
// IKE daemon listen on UDP 500 and 4500
const (
	maxQueueSize       int = 100000
	defaultIKEPort500  int = 500
	defaultIKEPort4500 int = 4500
	ChannelIDForPort500  = 1
	ChannelIDForPort4500 = 2
	)

type sendParameters struct {
	DstAddr *net.UDPAddr
	Length  int
	Payload []byte
}

var ikeLog *logrus.Entry

var sendChanToPort500 chan sendParameters // Chennel ID 1 (channelIDForPort500)
var mtxForChanPort500 sync.Mutex

var sendChanToPort4500 chan sendParameters // Chennel ID 2 (channelIDForPort4500)
var mtxForChanPort4500 sync.Mutex

func init() {
	// init logger
	ikeLog = logger.IKELog
	// init channel
	sendChanToPort500 = make(chan sendParameters, maxQueueSize)
	sendChanToPort4500 = make(chan sendParameters, maxQueueSize)
}

func configBindAddr(listenAddrPort500 *net.UDPAddr, listenAddrPort4500 *net.UDPAddr) {

	UeSelf := context.UE_Self()

	// Configure UDP port
	listenAddrPort500.Port, listenAddrPort4500.Port = defaultIKEPort500, defaultIKEPort4500

	// Configure IP address
	ip := net.ParseIP(UeSelf.IKEBindAddress)
	if ip != nil {
		ikeLog.Tracef("[IKE] Binding %v", ip)
		listenAddrPort500.IP, listenAddrPort4500.IP = ip, ip
	} else {
		ikeLog.Warn("[IKE] Invalid IKE bind IP address, binding 0.0.0.0")
		listenAddrPort500.IP, listenAddrPort4500.IP = net.IPv4zero, net.IPv4zero
	}
}

func Send(sendInfo *n3iwf_message.UDPSendInfoGroup, msg []byte) {
	if sendInfo.ChannelID == ChannelIDForPort500 {

		sendData := sendParameters{
			DstAddr: sendInfo.Addr,
			Length:  len(msg),
			Payload: msg,
		}

		mtxForChanPort500.Lock()
		sendChanToPort500 <- sendData
		mtxForChanPort500.Unlock()

	} else if sendInfo.ChannelID == ChannelIDForPort4500 {

		sendData := sendParameters{
			DstAddr: sendInfo.Addr,
			Length:  len(msg),
			Payload: msg,
		}

		mtxForChanPort4500.Lock()
		sendChanToPort4500 <- sendData
		mtxForChanPort4500.Unlock()

	} else {
		ikeLog.Error("[IKE] Send(): Invalid channel ID")
	}
}

func sender(channelID int, conn *net.UDPConn) {

	if channelID == ChannelIDForPort500 {
		for {

			sendData := <-sendChanToPort500

			n, err := conn.WriteToUDP(sendData.Payload, sendData.DstAddr)
			if err != nil {
				ikeLog.Errorf("[IKE] Sending data through UDP failed: %+v", err)
			}
			if n != sendData.Length {
				ikeLog.Warn("[IKE] There is data not being sent")
			}

		}
	} else if channelID == ChannelIDForPort4500 {
		for {

			sendData := <-sendChanToPort4500

			n, err := conn.WriteToUDP(sendData.Payload, sendData.DstAddr)
			if err != nil {
				ikeLog.Errorf("[IKE] Sending data through UDP failed: %+v", err)
			}
			if n != sendData.Length {
				ikeLog.Warn("[IKE] There is data not being sent")
			}

		}
	} else {
		ikeLog.Error("[IKE] sender(): Invalid channel ID")
	}
}

func reader(channelID int, conn *net.UDPConn) {

	if channelID > ChannelIDForPort4500 {
		ikeLog.Error("[IKE] Channel ID out of range")
		return
	}

	data := make([]byte, 65535)

	for {

		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			ikeLog.Errorf("[IKE] Read from UDP failed: %+v", err)
			continue
		}

		sendInfo := &n3iwf_message.UDPSendInfoGroup{
			ChannelID: channelID,
			Addr:      remoteAddr,
		}

		msg := n3iwf_message.HandlerMessage{
			Event:       n3iwf_message.EventN1UDPMessage,
			UDPSendInfo: sendInfo,
			Value:       data[:n],
		}

		n3iwf_message.SendMessage(msg)
	}

}

func Run() {

	listenAddrPort500 := new(net.UDPAddr)
	listenAddrPort4500 := new(net.UDPAddr)

	configBindAddr(listenAddrPort500, listenAddrPort4500)

	listenerPort500, err := net.ListenUDP("udp", listenAddrPort500)
	if err != nil {
		ikeLog.Errorf("[IKE] Listen on UDP socket failed: %+v", err)
		return
	}

	go reader(ChannelIDForPort500, listenerPort500)
	go sender(ChannelIDForPort500, listenerPort500)

	listenerPort4500, err := net.ListenUDP("udp", listenAddrPort4500)
	if err != nil {
		ikeLog.Errorf("[IKE] Listen on UDP socket failed: %+v", err)
		return
	}
	go reader(ChannelIDForPort4500, listenerPort4500)
	go sender(ChannelIDForPort4500, listenerPort4500)
}
*/
var ikeLog *logrus.Entry

func init() {
	// init logger
	ikeLog = logger.IKELog
}

func Run() error{

	//Resolver UDP for IKE SA INIT
	ip := context.UE_Self().IKEBindAddress
	udpAddrPort500, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:500", ip))
	if err != nil{
		ikeLog.Errorf("Ike resolver UDP port 500 address failed: %v", err)
		return errors.New("IKE service run failed")
	}

	udpAddrPort4500, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:4500", ip))
	if err != nil {
		ikeLog.Errorf("Ike resolver UDP port 500 address failed: %v", err)
	}

	var errChan chan error

	//Port 500
	errChan = make(chan error)
	go listenAndServe(udpAddrPort500, errChan)
	if err, ok := <-errChan; ok {
		ikeLog.Errorln(err)
		return errors.New("IKE service run failed")
	}

	//Port 4500
	errChan = make(chan error)
	go listenAndServe(udpAddrPort4500, errChan)
	if err, ok := <-errChan; ok {
		ikeLog.Errorln(err)
		return errors.New("Ike service run failed")
	}
	return nil
}

func listenAndServe(localAddr *net.UDPAddr, errChan chan<-error){
	listener, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		ikeLog.Errorf("Listen UDP failed: %+v", err)
		errChan <- errors.New("listenAndServe failed")
		return
	}

	close(errChan)

	data := make([] byte, 65535) //Why 65535

	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			ikeLog.Errorf("ReadFromUDP failed &+v", err)
			continue
		}

		forwardData := make([]byte, n)
		copy(forwardData, data[:n])

		go ike.Dispatch(listener, localAddr, remoteAddr, forwardData)
	}
}