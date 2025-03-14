package server

import (
	"context"
	"errors"
	"net"
	"strings"
)

func handle(resolver *net.Resolver, req *DNSMessage) (*DNSMessage, error) {
	flag := NewDNSFlag()
	flag.SetQR(1)
	flag.SetOpcode(req.hdr.flags.Opcode())
	flag.SetRD(req.hdr.flags.RD())
	if req.hdr.flags.Opcode() != 0 {
		flag.SetRCode(4)
	}

	var rrs []*ResourceRecord
	for _, q := range req.qs {
		if resolver == nil {
			rrs = append(rrs, &ResourceRecord{
				name:  q.labels,
				typ:   1,
				class: 1,
				ttl:   60,
				rdlen: 4,
				rdata: []byte{byte(8), byte(8), byte(8), byte(8)},
			})
			continue
		}

		ips, err := resolver.LookupIPAddr(context.Background(), strings.Join(q.labels, "."))
		if err != nil {
			return nil, errors.New("fail to lookup: " + err.Error())
		}

		for _, ip := range ips {
			rrs = append(rrs, &ResourceRecord{
				name:  q.labels,
				typ:   1,
				class: 1,
				ttl:   60,
				rdlen: 4,
				rdata: []byte(ip.IP),
			})
		}
	}

	m := &DNSMessage{
		hdr: &DNSHeader{
			id:      req.hdr.id,
			flags:   flag,
			qdcount: uint16(len(req.qs)),
			ancount: uint16(len(rrs)),
			nscount: 0,
			arcount: 0,
		},
		qs:  req.qs,
		rrs: rrs,
	}
	return m, nil
}
