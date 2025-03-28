package utils

func ToUint16(bs []byte) uint16 {
	return uint16(bs[0])<<8 | uint16(bs[1])
}

func ToUint32(bs []byte) uint32 {
	return uint32(bs[0])<<24 | uint32(bs[1])<<16 | uint32(bs[2])<<8 | uint32(bs[3])
}

func ToUint64(bs []byte) uint64 {
	return uint64(bs[0])<<56 | uint64(bs[1])<<48 | uint64(bs[2])<<40 | uint64(bs[3])<<32 | uint64(bs[4])<<24 | uint64(bs[5])<<16 | uint64(bs[6])<<8 | uint64(bs[7])
}

func FromUint16(n uint16) []byte {
	return []byte{byte(n >> 8), byte(n & 0xFF)}
}

func FromUint32(n uint32) []byte {
	return []byte{byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n & 0xFF)}
}
