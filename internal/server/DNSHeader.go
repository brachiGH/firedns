package server

import "github.com/brachiGH/firedns/internal/utils"

type DNSHeader struct {
	id      uint16
	flags   DNSFlag
	qdcount uint16 // Question count
	ancount uint16 // Answer count
	nscount uint16 // Authority count
	arcount uint16 // Additional count
}

type DNSFlag [2]byte

func NewDNSFlag() DNSFlag {
	return [2]byte{}
}

// Set Query/Response flag
func (f *DNSFlag) SetQR(v int) {
	if v != 0 {
		v = 1
	}
	f[0] = (f[0] & ^qrMask) | (byte(v) * qrMask)
}

// Getter for Opcode field
func (f *DNSFlag) Opcode() uint8 {
	return f[0] & opcodeMask >> 3
}

// Setter for Opcode field
func (f *DNSFlag) SetOpcode(code uint8) {
	code = code & 15
	f[0] = (f[0] & ^opcodeMask) | (code << 3)
}

// Getter for Recursion Desired field
func (f *DNSFlag) RD() int {
	v := f[0] & rdMask
	if v != 0 {
		return 1
	}
	return 0
}

// Setter for Recursion Desired field
func (f *DNSFlag) SetRD(v int) {
	if v != 0 {
		v = 1
	}
	f[0] = (f[0] & ^rdMask) | (byte(v) * rdMask)
}

// Getter for Recrsion Desired field
func (f *DNSFlag) SetRCode(code uint8) {
	code = code & 15
	f[1] = (f[1] & ^rcodeMask) | byte(code)
}

// Setter for Recursion Available field
func (f *DNSFlag) RA() int {
	v := f[1] & raMask
	if v != 0 {
		return 1
	}
	return 0
}

// Getter for Recursion Available field
func (f *DNSFlag) SetRA(v int) {
	if v != 0 {
		v = 1
	}
	f[1] = (f[1] & ^raMask) | (byte(v) * raMask)
}

// Getter for tauncated field
func (f *DNSFlag) TC() int {
	v := f[1] & tcMask
	if v != 0 {
		return 1
	}
	return 0
}

// Setter for Truncated field
func (f *DNSFlag) SetTC(v int) {
	if v != 0 {
		v = 1
	}
	f[1] = (f[1] & ^tcMask) | (byte(v) * tcMask)
}

func NewDNSHeader(data []byte) *DNSHeader {
	hdr := &DNSHeader{
		id:      utils.ToUint16(data[0:2]),
		qdcount: utils.ToUint16(data[4:6]),
		ancount: utils.ToUint16(data[6:8]),
		nscount: utils.ToUint16(data[8:12]),
		arcount: utils.ToUint16(data[10:12]),
	}
	copy(hdr.flags[:], data[2:4])
	return hdr
}
