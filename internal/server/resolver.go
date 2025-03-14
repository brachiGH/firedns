package server

import (
	"context"
	"net"
	"time"
)

func dial(ns string) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: time.Millisecond * time.Duration(10000),
		}
		return d.DialContext(ctx, network, ns)
	}
}

func newResolver(ns string) *net.Resolver {
	return &net.Resolver{
		Dial: dial(ns),
	}
}
