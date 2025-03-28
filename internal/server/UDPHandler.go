package server

import (
	"fmt"
	"net"

	"github.com/brachiGH/firedns/internal/transport"
	"github.com/brachiGH/firedns/internal/utils"
	"github.com/brachiGH/firedns/internal/utils/config"
	"github.com/brachiGH/firedns/monitor"
)

// checks for blocked domains and forwards the packet to dns resolver
func handle(data []byte, sourceIP net.IP) ([]byte, error) {
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
		return nil, fmt.Errorf("error parsing incoming message: %v", err)
	}

	lables := q.GetLables(data)
	blocked := CheckIfDomainIsBlocked(lables, sourceIP)

	if blocked {
		length := int(q.labelsEndPointer-q.labelsStartPointer) - 1
		if length > 0 {
			monitor.Droped(utils.ToIP(sourceIP), lables)

			dnsMessage := CreateBlockedDomainDNSMessage(data, hdr, q)
			return dnsMessage.AsBytes(), nil
		} else {
			return nil, fmt.Errorf("invalid QNAME")
		}
	}

	monitor.Passed(utils.ToIP(sourceIP), lables)

	arr := GetDNSAnswer__Cache(lables)
	if arr == nil {
		data, err = transport.ForwardPacketTo(data, config.UDP_ns_addr)
		if err != nil {
			return nil, fmt.Errorf("fail to lookup: %w", err)
		}

		arr = DNSAnswerFromBytes(data, q)
		if arr != nil {
			go AddDNSAnswer__Cache(lables, arr)
		}
	} else {
		dnsMessage := ReuseDNSAnswer__Cache(data, hdr, q, arr)
		return dnsMessage.AsBytes(), nil
	}

	return data, nil
}
