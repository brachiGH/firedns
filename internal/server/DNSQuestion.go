package server

import "github.com/brachiGH/firedns/internal/utils"

type ResourceRecord struct {
	name  []string
	typ   uint16
	class uint16
	ttl   uint32
	rdlen uint16
	rdata []byte
}

type DNSQuestion struct {
	labels []string
	typ    uint16
	class  uint16
}

func (r *ResourceRecord) AsBytes() []byte {
	var bs []byte
	for _, l := range r.name {
		bs = append(bs, byte(len(l)))
		bs = append(bs, []byte(l)...)
	}
	bs = append(bs, byte(0))
	bs = append(bs, utils.FromUint16(r.typ)...)
	bs = append(bs, utils.FromUint16(r.class)...)
	bs = append(bs, utils.FromUint32(r.ttl)...)
	bs = append(bs, utils.FromUint16(r.rdlen)...)
	bs = append(bs, r.rdata...)
	return bs
}

func (q *DNSQuestion) AsBytes() []byte {
	var bs []byte
	for _, l := range q.labels {
		bs = append(bs, byte(len(l)))
		bs = append(bs, []byte(l)...)
	}
	bs = append(bs, byte(0))
	bs = append(bs, utils.FromUint16(q.typ)...)
	bs = append(bs, utils.FromUint16(q.class)...)
	return bs
}

func NewDNSQuestion(start int, data []byte) (*DNSQuestion, int) {
	i := start
	labels := []string{}
	pointerEnd := -1
	for {
		l := int(data[i])
		i++
		if l == 0 {
			if pointerEnd == -1 {
				return &DNSQuestion{
					labels: labels,
					typ:    utils.ToUint16(data[i : i+2]),
					class:  utils.ToUint16(data[i+2 : i+4]),
				}, i + 4
			} else {
				return &DNSQuestion{
					labels: labels,
					typ:    utils.ToUint16(data[pointerEnd : pointerEnd+2]),
					class:  utils.ToUint16(data[pointerEnd+2 : pointerEnd+4]),
				}, pointerEnd + 4
			}
		}
		// handling DNS label compression, as described in RFC 1035, Section 4.1.4.
		if (l >> 6) == 3 {
			pointerEnd = i + 1
			i = (l & ((1 << 6) - 1) << 8) | int(data[i])
			continue
		}
		labels = append(labels, string(data[i:i+l]))
		i = i + l
	}
}
