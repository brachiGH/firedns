package server

import (
	"github.com/brachiGH/firedns/internal/utils"
)

type DNSAnswer struct {
	name  []byte
	typ   uint16
	class uint16
	ttl   uint32
	rdlen uint16
	rdata []byte
}

func DNSAnswerFromBytes(data []byte, q *DNSQuestion) *DNSAnswer {
	answerCount := int(data[6])<<8 | int(data[7])
	if answerCount != 1 {
		return nil
	}

	arrStartPointer := q.labelsEndPointer + QueryInfomationBytesLength + 1
	name := []byte{0xC0, 0x0C}
	if data[arrStartPointer] != 0xC0 {
		// note: not test but it should work fine
		name = data[q.labelsStartPointer : q.labelsEndPointer+1]
		arrStartPointer += q.labelsStartPointer - q.labelsEndPointer
	} else {
		arrStartPointer += 2
	}

	if uint8(len(data)) < arrStartPointer+11 {
		return nil
	}

	return &DNSAnswer{
		name:  name,
		typ:   utils.ToUint16(data[arrStartPointer : arrStartPointer+2]),
		class: utils.ToUint16(data[arrStartPointer+2 : arrStartPointer+4]),
		ttl:   utils.ToUint32(data[arrStartPointer+4 : arrStartPointer+8]),
		rdlen: utils.ToUint16(data[arrStartPointer+8 : arrStartPointer+10]),
		rdata: data[arrStartPointer+10:],
	}
}

func (r *DNSAnswer) AsBytes() []byte {
	var bs []byte
	bs = append(bs, r.name...)
	bs = append(bs, utils.FromUint16(r.typ)...)
	bs = append(bs, utils.FromUint16(r.class)...)
	bs = append(bs, utils.FromUint32(r.ttl)...)
	bs = append(bs, utils.FromUint16(r.rdlen)...)
	bs = append(bs, r.rdata...)
	return bs
}
