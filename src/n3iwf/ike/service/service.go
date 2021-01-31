package service

import (
	"errors"
	"net"

	"free5gc/src/n3iwf/context"
	"free5gc/src/n3iwf/ike"
	"free5gc/src/n3iwf/logger"

	"github.com/sirupsen/logrus"
)

var ikeLog *logrus.Entry

func init() {
	// init logger
	ikeLog = logger.IKELog
}

func Run() error {
	// Resolve UDP addresses
	ip := context.N3IWFSelf().IKEBindAddress
	udpAddrPort500, err := net.ResolveUDPAddr("udp", ip+":500")
	if err != nil {
		ikeLog.Errorf("Resolve UDP address failed: %+v", err)
		return errors.New("IKE service run failed")
	}
	udpAddrPort4500, err := net.ResolveUDPAddr("udp", ip+":4500")
	if err != nil {
		ikeLog.Errorf("Resolve UDP address failed: %+v", err)
		return errors.New("IKE service run failed")
	}

	// Listen and serve
	var errChan chan error

	// Port 500
	errChan = make(chan error)
	go listenAndServe(udpAddrPort500, errChan)
	if err, ok := <-errChan; ok {
		ikeLog.Errorln(err)
		return errors.New("IKE service run failed")
	}

	// Port 4500
	errChan = make(chan error)
	go listenAndServe(udpAddrPort4500, errChan)
	if err, ok := <-errChan; ok {
		ikeLog.Errorln(err)
		return errors.New("IKE service run failed")
	}

	return nil
}

func listenAndServe(localAddr *net.UDPAddr, errChan chan<- error) {
	listener, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		ikeLog.Errorf("Listen UDP failed: %+v", err)
		errChan <- errors.New("listenAndServe failed")
		return
	}

	close(errChan)

	data := make([]byte, 65535)

	for {

		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			ikeLog.Errorf("ReadFromUDP failed: %+v", err)
			continue
		}

		forwardData := make([]byte, n)
		copy(forwardData, data[:n])

		go ike.Dispatch(listener, localAddr, remoteAddr, forwardData)

	}
}
