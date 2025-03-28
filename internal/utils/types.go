package utils

import (
	"net"
)

// Label represents a domain name label.
type Lable string

type IP uint64

func LablesToStrings(lables []Lable) []string {
	result := make([]string, len(lables))
	for i, label := range lables {
		result[i] = string(label)
	}
	return result
}

func ToIP(bs net.IP) IP {
	bs4 := bs.To4()
	var ip uint64
	if bs4 == nil {
		bs6 := bs.To16()
		ip = uint64(bs6[7])<<56 | uint64(bs6[6])<<48 | uint64(bs6[5])<<40 | uint64(bs6[4])<<32 | uint64(bs6[3])<<24 | uint64(bs6[2])<<16 | uint64(bs6[1])<<8 | uint64(bs6[0])
	} else {
		ip = uint64(bs4[3])<<24 | uint64(bs4[2])<<16 | uint64(bs4[1])<<8 | uint64(bs4[0])
	}
	return IP(ip)
}
