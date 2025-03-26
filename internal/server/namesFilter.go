package server

func CheckIfDomainIsBlocked(qs *DNSQuestion, data []byte) bool {

	return false
}

func CreateBlockedDomainDNSMessage(data []byte, hdr *DNSHeader, q *DNSQuestion) *DNSMessage {
	arr := &AnswerResourceRecord{
		name:  data[q.labelsStartPointer:q.labelsEndPointer],
		typ:   1,
		class: 1,
		ttl:   60,
		rdlen: 4,
		rdata: []byte("0.0.0.0"),
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
		q:            q,
		AnswerRR:     arr,
		AdditionalRR: UDP_Response_Additional_Records,
	}

	return dnsMessage
}
