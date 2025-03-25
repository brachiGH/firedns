package server

import (
	"errors"
)

// parsing dns header and Questions
// returns header content and list raw lables list
func NewDNSMessage(data []byte) ([]*DNSQuestion, error) {
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
	for i := 0; i < int(hdr.qdcount) && qoffset <= maxPacketsize; i++ {
		qs[i], qoffset = NewDNSQuestion(qoffset, data)
	}

	return qs, nil
}
