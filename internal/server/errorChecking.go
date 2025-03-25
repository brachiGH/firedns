package server

import (
	"errors"
)

func checkHeaderErrors(hdr *DNSHeader) error {
	opcode := hdr.flags.Opcode()
	if opcode == 1 {
		//RFC 3425 has obsoleted IQUERY due to the burden it placed on servers and security concerns (RFC 1035 defined inverse queries (IQUERY with OPCODE 1) )
		return errors.New("obsoleted query no longer suppored")
	}
	if opcode == 0 && hdr.qdcount != 1 {
		//RFC 9619 Update: This document explicitly states that a DNS message with OPCODE = 0 (QUERY) MUST NOT include a QDCOUNT parameter whose value is diffrent than 1
		return errors.New("unvalide qd_count")
	}

	return nil
}
