package utils

func ToUint16(bs []byte) uint16 {
	return uint16(bs[0])<<8 | uint16(bs[1])
}

func FromUint16(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n & 255)}
}

func FromUint32(n uint32) []byte {
	return []byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n & 255)}
}
