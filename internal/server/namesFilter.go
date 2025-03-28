package server

import (
	"net"

	"github.com/brachiGH/firedns/internal/utils"
	"github.com/brachiGH/firedns/internal/utils/config"
)

func CheckIfDomainIsBlocked(lables []utils.Lable, IP net.IP) bool {

	return false
}

func CreateBlockedDomainDNSMessage(data []byte, hdr *DNSHeader, q *DNSQuestion) *DNSMessage {
	arr := &DNSAnswer{
		name:  []byte{0xC0, 0x0C},
		typ:   1,
		class: 1,
		ttl:   60,
		rdlen: 4,
		rdata: []byte{0x00, 0x00, 0x00, 0x00},
	}

	hdr.ancount = 1
	hdr.arcount = 1

	rd := hdr.flags.RD()

	hdr.flags = NewDNSFlag()
	hdr.flags.SetQR(1)
	hdr.flags.SetRD(rd)
	hdr.flags.SetRA(1)

	dnsMessage := &DNSMessage{
		hdr:          hdr,
		q:            data[q.labelsStartPointer : q.labelsEndPointer+1+QueryInfomationBytesLength],
		AnswerRR:     arr,
		AdditionalRR: config.UDP_Response_Additional_Records,
	}

	return dnsMessage
}
