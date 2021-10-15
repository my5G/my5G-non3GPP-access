package service

import (
	"errors"
	"io"
	"time"

	"git.cs.nctu.edu.tw/calee/sctp"
	"github.com/sirupsen/logrus"

	"free5gc/src/n3iwf/context"
	"free5gc/src/n3iwf/logger"
	"free5gc/src/n3iwf/ngap"
	"free5gc/src/n3iwf/ngap/message"
)

var ngapLog *logrus.Entry

const NGAP_PPID_BigEndian = 0x3c000000

func init() {
	ngapLog = logger.NgapLog
}

// Run start the N3IWF SCTP process.
func Run() error {
	// n3iwf context
	n3iwfSelf := context.N3IWFSelf()
	// load amf SCTP address slice
	amfSCTPAddresses := n3iwfSelf.AMFSCTPAddresses

	localAddr := new(sctp.SCTPAddr)

	for _, remoteAddr := range amfSCTPAddresses {
		errChan := make(chan error)
		go listenAndServe(localAddr, remoteAddr, errChan)
		if err, ok := <-errChan; ok {
			ngapLog.Errorln(err)
			return errors.New("NGAP service run failed")
		}
	}

	return nil
}

func listenAndServe(localAddr, remoteAddr *sctp.SCTPAddr, errChan chan<- error) {
	var conn *sctp.SCTPConn
	var err error

	// Connect the session
	for i := 0; i < 3; i++ {
		conn, err = sctp.DialSCTP("sctp", localAddr, remoteAddr)
		if err != nil {
			ngapLog.Errorf("[SCTP] DialSCTP(): %+v", err)
		} else {
			break
		}

		if i != 2 {
			ngapLog.Info("Retry to connect AMF after 1 second...")
			time.Sleep(1 * time.Second)
		} else {
			ngapLog.Debugf("[SCTP] AMF SCTP address: %+v", remoteAddr.String())
			errChan <- errors.New("Failed to connect to AMF.")
			return
		}
	}

	// Set default sender SCTP infomation sinfo_ppid = NGAP_PPID = 60
	info, err := conn.GetDefaultSentParam()
	if err != nil {
		ngapLog.Errorf("[SCTP] GetDefaultSentParam(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			ngapLog.Errorf("conn close error in GetDefaultSentParam(): %+v", errConn)
		}
		errChan <- errors.New("Get socket infomation failed.")
		return
	}
	info.PPID = NGAP_PPID_BigEndian
	err = conn.SetDefaultSentParam(info)
	if err != nil {
		ngapLog.Errorf("[SCTP] SetDefaultSentParam(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			ngapLog.Errorf("conn close error in SetDefaultSentParam(): %+v", errConn)
		}
		errChan <- errors.New("Set socket parameter failed.")
		return
	}

	// Subscribe receiver SCTP information
	err = conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		ngapLog.Errorf("[SCTP] SubscribeEvents(): %+v", err)
		errConn := conn.Close()
		if errConn != nil {
			ngapLog.Errorf("conn close error in SubscribeEvents(): %+v", errConn)
		}
		errChan <- errors.New("Subscribe SCTP event failed.")
		return
	}

	// Send NG setup request
	go message.SendNGSetupRequest(conn)

	close(errChan)

	data := make([]byte, 65535)

	for {
		n, info, err := conn.SCTPRead(data)

		if err != nil {
			ngapLog.Debugf("[SCTP] AMF SCTP address: %+v", conn.RemoteAddr().String())
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				ngapLog.Warn("[SCTP] Close connection.")
				errConn := conn.Close()
				if errConn != nil {
					ngapLog.Errorf("conn close error: %+v", errConn)
				}
				return
			}
			ngapLog.Errorf("[SCTP] Read from SCTP connection failed: %+v", err)
		} else {
			ngapLog.Tracef("[SCTP] Successfully read %d bytes.", n)

			if info == nil || info.PPID != NGAP_PPID_BigEndian {
				ngapLog.Warn("Received SCTP PPID != 60")
				continue
			}

			forwardData := make([]byte, n)
			copy(forwardData, data[:n])

			go ngap.Dispatch(conn, forwardData)
		}
	}
}
