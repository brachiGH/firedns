package utils

func ToUint16(bs []byte) uint16 {
	return uint16(bs[0])<<8 | uint16(bs[1])
}

func ToUint32(bs []byte) uint32 {
	return uint32(bs[0])<<24 | uint32(bs[1])<<16 | uint32(bs[2])<<8 | uint32(bs[3])
}

func FromUint16(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n & 0xFF)}
}

func FromUint32(n uint32) []byte {
	return []byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n & 0xFF)}
}
