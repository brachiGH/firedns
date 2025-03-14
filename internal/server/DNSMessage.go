package server

import (
	"errors"

	"github.com/brachiGH/firedns/internal/utils"
)

type DNSMessage struct {
	hdr *DNSHeader
	qs  []*DNSQuestion
	rrs []*ResourceRecord
}

func (m *DNSMessage) AsBytes() []byte {
	b := make([]byte, hdrSize)
	copy(b[0:2], utils.FromUint16(m.hdr.id))
	copy(b[2:4], m.hdr.flags[:])
	copy(b[4:6], utils.FromUint16(m.hdr.qdcount))
	copy(b[6:8], utils.FromUint16(m.hdr.ancount))
	copy(b[8:10], utils.FromUint16(m.hdr.nscount))
	copy(b[10:12], utils.FromUint16(m.hdr.arcount))

	for _, q := range m.qs {
		b = append(b, q.AsBytes()...)
	}
	for _, r := range m.rrs {
		b = append(b, r.AsBytes()...)
	}
	return b
}

func NewDNSMessage(data []byte) (*DNSMessage, error) {
	if len(data) < hdrSize {
		return nil, errors.New("message is too short")
	}

	hdr := NewDNSHeader(data)

	err := checkHeaderErrors(hdr)
	if err != nil {
		return nil, err
	}

	qoffset := hdrSize
	qs := make([]*DNSQuestion, hdr.qdcount)
	for i := 0; i < int(hdr.qdcount); i++ {
		qs[i], qoffset = NewDNSQuestion(qoffset, data)
	}

	return &DNSMessage{
		hdr: hdr,
		qs:  qs,
	}, nil
}
