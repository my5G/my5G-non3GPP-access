package pfcp

import (
	"fmt"
	"free5gc/lib/tlv"
)

func (m *Message) Marshal() ([]byte, error) {
	var headerBytes []byte
	var bodyBytes []byte

	var err error

	if bodyBytes, err = tlv.Marshal(m.Body); err != nil {
		return nil, fmt.Errorf("pfcp: marshal message type %d body failed: %s", m.Header.MessageType, err)
	}

	if m.Header.S&1 != 0 {
		// 8 (SEID) + 3 (Sequence Number) + 1 (Message Priority and Spare)
		m.Header.MessageLength = 12
	} else {
		// 3 (Sequence Number) + 1 (Message Priority and Spare)
		m.Header.MessageLength = 4
	}

	m.Header.MessageLength += uint16(len(bodyBytes))
	if headerBytes, err = m.Header.MarshalBinary(); err != nil {
		return nil, fmt.Errorf("pfcp: marshal message type %d header failed: %s", m.Header.MessageType, err)
	}

	return append(headerBytes, bodyBytes[:]...), nil
}
