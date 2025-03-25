package server

// this struct stores the pointer to the QNAME in each question
// note: a QNAME is represented as a sequence of length-delimited labels
type DNSQuestion struct {
	labelsStartPointer uint8
	labelsEndPointer   uint8
}

func NewDNSQuestion(start int, data []byte) (*DNSQuestion, int) {
	i := start
	pointerEnd := -1
	for {
		l := int(data[i])
		i++
		if l == 0 {
			var result int
			if pointerEnd == -1 {
				result = i + 4
			} else {
				result = pointerEnd + 4
			}

			return &DNSQuestion{
				labelsStartPointer: uint8(start),
				labelsEndPointer:   uint8(i - 1),
			}, result
		}

		// handling DNS label compression, as described in RFC 1035, Section 4.1.4.
		if (l >> 6) == 3 {
			pointerEnd = i + 1
			i = (l & ((1 << 6) - 1) << 8) | int(data[i])
			continue
		}
		i = i + l
	}
}
