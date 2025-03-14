package server

const hdrSize = 12

const (
	qroffset   int  = 7
	qrMask     byte = 1 << qroffset
	rdoffset   int  = 0
	rdMask     byte = 1 << rdoffset
	opcodeMask byte = ((1 << 4) - 1) << 3
	rcodeMask  byte = (1 << 4) - 1
)
