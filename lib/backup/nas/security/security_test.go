package security_test

import (
	"encoding/binary"
	"free5gc/lib/CommonConsumerTestData/AMF/TestAmf/TestNeaNia"
	"free5gc/lib/nas/security"
	"testing"
)

func TestNEA1(t *testing.T) {
	for _, ts := range TestNeaNia.TestNEA1Table {
		obs, _ := security.NEA1(ts.Ck, ts.CountC, ts.Bearer, ts.Direction, ts.Ibs, ts.Length)
		for i, v := range obs {
			if v != ts.Obs[i] {
				t.Errorf("%#x != %#x\n", v, ts.Obs[i])
			}
		}
	}
}

func TestNIA1(t *testing.T) {
	for _, ts := range TestNeaNia.TestNIA1Table {
		buf, _ := security.NIA1(ts.Ik, ts.CountI, ts.Bearer, ts.Direction, ts.Msg, ts.Length)
		mac := binary.BigEndian.Uint32(buf)
		if mac != ts.MacI {
			t.Errorf("%#x != %#x\n", mac, ts.MacI)
		}
	}
}
