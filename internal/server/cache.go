package server

import (
	"sync"
	"time"

	"github.com/brachiGH/firedns/internal/utils"
)

type nameSpaceNode struct {
	arr      DNSAnswer
	children map[utils.Lable]*nameSpaceNode
}

func NewNameSpaceNode(arr DNSAnswer) *nameSpaceNode {
	return &nameSpaceNode{arr: arr, children: make(map[utils.Lable]*nameSpaceNode)}
}

var nameSpaceRoot = NewNameSpaceNode(DNSAnswer{})
var namesSpaceWG sync.WaitGroup

func GetDNSAnswer__Cache(lables []utils.Lable) (arr *DNSAnswer) {
	namesSpaceWG.Wait() // block until any ongoing cache clearing operation completes

	parent := nameSpaceRoot
	for i := len(lables) - 1; i >= 0; i-- {
		parent = parent.children[lables[i]]

		if parent == nil {
			return nil
		}
	}

	return &parent.arr
}

func AddDNSAnswer__Cache(lables []utils.Lable, arr *DNSAnswer) {
	namesSpaceWG.Wait() // block until any ongoing cache clearing operation completes

	parent := nameSpaceRoot

	for i := len(lables) - 1; i >= 0; i-- {
		child, exists := parent.children[lables[i]]
		if !exists {
			child = NewNameSpaceNode(DNSAnswer{})
			parent.children[lables[i]] = child
		}
		parent = child
	}

	parent.arr = *arr
}

func ReuseDNSAnswer__Cache(data []byte, hdr *DNSHeader, q *DNSQuestion, arr *DNSAnswer) *DNSMessage {
	hdr.ancount = 1

	rd := hdr.flags.RD()

	hdr.flags = NewDNSFlag()
	hdr.flags.SetQR(1)
	hdr.flags.SetRD(rd)
	hdr.flags.SetRA(1)

	dnsMessage := &DNSMessage{
		hdr:      hdr,
		q:        data[q.labelsStartPointer : q.labelsEndPointer+1+QueryInfomationBytesLength],
		AnswerRR: arr,
	}

	return dnsMessage
}

func ClearCache_Routine() {
	tick := time.Tick(time.Minute)
	for range tick {
		namesSpaceWG.Add(1)
		nameSpaceRoot.children = make(map[utils.Lable]*nameSpaceNode)
		namesSpaceWG.Done()
	}
}
