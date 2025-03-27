package server

import "strings"

// this struct stores the pointer to the QNAME in each question
// note: a QNAME is represented as a sequence of length-delimited labels
type DNSQuestion struct {
	labelsStartPointer uint8
	labelsEndPointer   uint8 // Note: this pointer dosn't pointer to the null byte at the end of the QNAME
}

// Label represents a domain name label.
type Lable string

// NewDNSQuestion parses a QNAME stored in data starting at index start.
// It returns a pointer to DNSQuestion and the index immediately after the question.
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

// GetLabels returns a slice of labels (strings) for the QNAME based on the DNSQuestion pointers.
func (q *DNSQuestion) GetLables(data []byte) []Lable {
	labels := make([]Lable, 0)

	for i := q.labelsStartPointer; i < q.labelsEndPointer; i++ {
		lableSize := data[i]

		lable := strings.ToLower(string(data[i+1 : i+1+lableSize]))
		labels = append(labels, Lable(lable))

		i += lableSize
	}

	return labels
}
