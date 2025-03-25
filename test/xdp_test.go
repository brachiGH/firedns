package test

import (
	"testing"

	"github.com/brachiGH/firedns/monitor"
)

func TestXDP(t *testing.T) {
	var xdp monitor.XDPobj
	err := xdp.Load()
	if err != nil {
		t.Fatal("Failed to load and link XDP: ", err)
	}
	err = xdp.Link()
	if err != nil {
		t.Fatal("Failed to link XDP: ", err)
	}
	defer xdp.UnloadAndCLoseLink()
}
