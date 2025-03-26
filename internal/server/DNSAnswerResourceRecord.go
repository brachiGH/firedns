package server

import (
	"github.com/brachiGH/firedns/internal/utils"
)

type AnswerResourceRecord struct {
	name  []byte
	typ   uint16
	class uint16
	ttl   uint32
	rdlen uint16
	rdata []byte
}

func (r *AnswerResourceRecord) AsBytes() []byte {
	var bs []byte
	bs = append(bs, byte(0))
	bs = append(bs, r.name...)
	bs = append(bs, utils.FromUint16(r.typ)...)
	bs = append(bs, utils.FromUint16(r.class)...)
	bs = append(bs, utils.FromUint32(r.ttl)...)
	bs = append(bs, utils.FromUint16(r.rdlen)...)
	bs = append(bs, r.rdata...)
	return bs
}
