package server

import (
	"errors"

	"github.com/brachiGH/firedns/internal/utils"
)

type DNSMessage struct {
	hdr          *DNSHeader
	q            []byte
	AnswerRR     *DNSAnswer
	AdditionalRR []byte
}

// parsing dns header and Questions
// returns header content and list raw lables list
func NewDNSMessage(data []byte) (*DNSQuestion, error) {
	qoffset := HdrSize
	/** RFC 9619 Update (2024-7-24): This document explicitly states that a DNS message with OPCODE = 0 (QUERY)
	MUST NOT include a QDCOUNT parameter whose value is diffrent than 1

	// qs := make([]*DNSQuestion, hdr.qdcount)
	// for i := 0; i < int(hdr.qdcount) && qoffset <= maxPacketsize; i++ {
	// 	qs[i], qoffset = NewDNSQuestion(qoffset, data)
	// }
	**/
	var qs *DNSQuestion
	qs, qoffset = NewDNSQuestion(qoffset, data)

	if qs == nil || qoffset > len(data) {
		return nil, errors.New("invalid question")
	}

	return qs, nil
}

func (m *DNSMessage) AsBytes() []byte {
	b := make([]byte, HdrSize)
	copy(b[0:2], utils.FromUint16(m.hdr.id))
	copy(b[2:4], m.hdr.flags[:])
	copy(b[4:6], utils.FromUint16(m.hdr.qdcount))
	copy(b[6:8], utils.FromUint16(m.hdr.ancount))
	copy(b[8:10], utils.FromUint16(m.hdr.nscount))
	copy(b[10:12], utils.FromUint16(m.hdr.arcount))
	b = append(b, m.q...)
	b = append(b, m.AnswerRR.AsBytes()...)
	b = append(b, m.AdditionalRR...)
	return b
}
