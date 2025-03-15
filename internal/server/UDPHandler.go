package server

import (
	"bytes"
	"errors"

	"github.com/brachiGH/firedns/internal/transport"
)

// checks for blocked domains and forwards the packet to dns resolver
func handle(data []byte, qs []*DNSQuestion) ([]byte, error) {

	for _, q := range qs {
		blocked := CheckIfDomainIsBlocked(q, data)

		if blocked {
			length := int(q.labelsEndPointer-q.labelsStartPointer) - 1
			if length > 0 {
				copy(data[q.labelsStartPointer:], append([]byte{byte(length)}, bytes.Repeat([]byte{'x'}, length)...))
			} else {
				return []byte{}, errors.New("invalid QNAME")
			}
		}
	}

	data, err := transport.ForwardPacketTo(data, upd_ns)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
