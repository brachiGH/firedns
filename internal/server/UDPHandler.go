package server

import (
	"fmt"

	"github.com/brachiGH/firedns/internal/transport"
)

// checks for blocked domains and forwards the packet to dns resolver
func handle(data []byte) ([]byte, error) {
	if len(data) < hdrSize {
		return nil, fmt.Errorf("message is too short")
	}

	hdr := NewDNSHeader(data)

	err := checkHeaderErrors(hdr)
	if err != nil {
		return nil, err
	}

	q, err := NewDNSMessage(data)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing incoming message: %v", err)
	}

	blocked := CheckIfDomainIsBlocked(q, data)

	if blocked {
		length := int(q.labelsEndPointer-q.labelsStartPointer) - 1
		if length > 0 {
			dnsMessage := CreateBlockedDomainDNSMessage(data, hdr, q)
			return dnsMessage.AsBytes(data), nil
		} else {
			return []byte{}, fmt.Errorf("invalid QNAME")
		}
	}

	data, err = transport.ForwardPacketTo(data, UDP_ns_addr)
	if err != nil {
		return []byte{}, fmt.Errorf("fail to lookup: %w", err)
	}

	// var rrs []*ResourceRecord
	// for _, q := range qs {
	// 	label := string(data[q.labelsStartPointer:q.labelsEndPointer])
	// rrs = append(rrs, &ResourceRecord{
	// 	name:  label,
	// 	typ:   1,
	// 	class: 1,
	// 	ttl:   60,
	// 	rdlen: 4,
	// 	rdata: []byte(ip.IP),
	// })
	// }

	return data, nil
}
