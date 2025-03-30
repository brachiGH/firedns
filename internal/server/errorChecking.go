package server

import "errors"

func checkHeaderErrors(hdr *DNSHeader) error {
	opcode := hdr.flags.Opcode()
	if hdr.qdcount != 1 {
		//RFC 9619 Update: This document explicitly states that a DNS message with OPCODE = 0 (QUERY) MUST NOT include a QDCOUNT parameter whose value is diffrent than 1
		// FireDNS only supports qdcount with question equal to 1
		return errors.New("invalid query count")
	}

	if opcode == 1 {
		//RFC 3425 has obsoleted IQUERY due to the burden it placed on servers and security concerns (RFC 1035 defined inverse queries (IQUERY with OPCODE 1) )
		return errors.New("obsoleted query no longer suppored")
	}

	if opcode == 2 {
		//This value indicates a server status request. It was listed as an option in [RFC1034]. The sources do not provide further details on its current status or usage.
		return errors.New("server status request not supported")
	}

	if ((opcode < 5) && (opcode > 2)) || (opcode > 6) {
		// 3-4: These values were reserved for future use according to [RFC1035].
		// 7-15: These values were also reserved for future use according to [RFC1035].
		return errors.New("reserved opcode")
	}

	if opcode == 5 {
		//This opcode identifies a DNS Update message. DNS Update is used to dynamically update resource records in an authoritative zone.
		return errors.New("update query not supported")
	}

	if hdr.flags.QR() != 0 {
		return errors.New("invalid dns message")
	}

	return nil
}
