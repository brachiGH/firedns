package server

const (
	HdrSize                    = 12
	QueryInfomationBytesLength = 4 // 2 bytes for QTYPE and 2 bytes for QCLASS

	qroffset int  = 7 // Set Query/Response flag
	qrMask   byte = 1 << qroffset

	tcoffset int  = 2 // Truncated flag
	tcMask   byte = 1 << tcoffset

	rdoffset int  = 0 // Recursion Desired flag
	rdMask   byte = 1 << rdoffset

	raoffset int  = 7 // Recursion Available flag
	raMask   byte = 1 << raoffset

	opcodeMask byte = ((1 << 4) - 1) << 3 // Opcode field
	rcodeMask  byte = (1 << 4) - 1

	MaxPacketsize int = 512
)
