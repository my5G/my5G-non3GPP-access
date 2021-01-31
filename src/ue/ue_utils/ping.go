package ue_utils

import (
	"github.com/sparrc/go-ping"
	"time"
	"github.com/sirupsen/logrus"
	"free5gc/src/ue/logger"
)

var log *logrus.Entry

func init(){
	log = logger.AppLog
}

func Ping(srcAddr, destAddr string){
	// Ping remote
	pinger, err := ping.NewPinger(destAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Run with root
	pinger.SetPrivileged(true)

	pinger.OnRecv = func(pkt *ping.Packet) {
		log.Infof("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		log.Infof("\n--- %s ping statistics ---\n", stats.Addr)
		log.Infof("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		log.Infof("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	pinger.Count = 5
	pinger.Timeout = 10 * time.Second
	pinger.Source = "60.60.0.1"

	time.Sleep(3 * time.Second)

	pinger.Run()

	time.Sleep(1 * time.Second)

	stats := pinger.Statistics()
	if stats.PacketsSent != stats.PacketsRecv {
		log.Fatal("Ping Failed")
	}else{
		log.Infoln("Ping Succeed")
	}
}
