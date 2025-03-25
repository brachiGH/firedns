package server

import "github.com/brachiGH/firedns/internal/utils"

type DNSHeader struct {
	id      uint16
	flags   DNSFlag
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
}

type DNSFlag [2]byte

func NewDNSFlag() DNSFlag {
	return [2]byte{}
}

func (f *DNSFlag) SetQR(v int) {
	if v != 0 {
		v = 1
	}
	f[0] = (f[0] & ^qrMask) | (byte(v) * qrMask)
}

func (f *DNSFlag) Opcode() uint8 {
	return f[0] & opcodeMask >> 3
}

func (f *DNSFlag) SetOpcode(code uint8) {
	code = code & 15
	f[0] = (f[0] & ^opcodeMask) | (code << 3)
}

func (f *DNSFlag) RD() int {
	v := f[0] & rdMask
	if v != 0 {
		return 1
	}
	return 0
}

func (f *DNSFlag) SetRD(v int) {
	if v != 0 {
		v = 1
	}
	f[0] = (f[0] & ^rdMask) | (byte(v) * rdMask)
}
func (f *DNSFlag) SetRCode(code uint8) {
	code = code & 15
	f[1] = (f[1] & ^rcodeMask) | byte(code)
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
