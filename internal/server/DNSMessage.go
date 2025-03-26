package server

import (
	"github.com/brachiGH/firedns/internal/utils"
)

type DNSMessage struct {
	hdr          *DNSHeader
	q            *DNSQuestion
	AnswerRR     *AnswerResourceRecord
	AdditionalRR []byte
}

// parsing dns header and Questions
// returns header content and list raw lables list
func NewDNSMessage(data []byte) (*DNSQuestion, error) {
	qoffset := hdrSize
	/** RFC 9619 Update (2024-7-24): This document explicitly states that a DNS message with OPCODE = 0 (QUERY)
	MUST NOT include a QDCOUNT parameter whose value is diffrent than 1

	// qs := make([]*DNSQuestion, hdr.qdcount)
	// for i := 0; i < int(hdr.qdcount) && qoffset <= maxPacketsize; i++ {
	// 	qs[i], qoffset = NewDNSQuestion(qoffset, data)
	// }
	**/
	var qs *DNSQuestion
	qs, _ = NewDNSQuestion(qoffset, data)

	return qs, nil
}

func (m *DNSMessage) AsBytes(data []byte) []byte {
	b := make([]byte, hdrSize)
	copy(b[0:2], utils.FromUint16(m.hdr.id))
	copy(b[2:4], m.hdr.flags[:])
	copy(b[4:6], utils.FromUint16(m.hdr.qdcount))
	copy(b[6:8], utils.FromUint16(m.hdr.ancount))
	copy(b[8:10], utils.FromUint16(m.hdr.nscount))
	copy(b[10:12], utils.FromUint16(m.hdr.arcount))
	b = append(b, data[m.q.labelsStartPointer:m.q.labelsEndPointer+1+QueryInfomationBytesLength]...)
	b = append(b, m.AnswerRR.AsBytes()...)
	b = append(b, m.AdditionalRR...)
	return b
}
