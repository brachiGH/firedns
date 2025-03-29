package test

import (
	"net"
	"testing"
	"time"

	"github.com/brachiGH/firedns/internal/transport"
	"github.com/brachiGH/firedns/internal/utils/config"
)

func TestDNSudpMessage(t *testing.T) {
	go transport.Upd_dns_server()

	// Wait for the UDP server to start
	time.Sleep(1 * time.Second)

	// Connect to a DoT server (firedns's local on port 2053)
	addr, err := net.ResolveUDPAddr("udp", config.UDP_local_ns_addr)
	if err != nil {
		t.Fatal("Failed to error resolving UDP address: ", err)
	}

	// Create UDP Connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Fatal("error connecting to UDP server: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Log("error closing connection: ", err)
		}
	}()

	// Minimal DNS Query (Query for google.com, type A)
	data := []byte{
		0x00, 0x01, // Transaction ID
		0x00, 0x00, // Flags
		0x00, 0x01, // Question Count
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x3, 0x63, 0x6f, 0x6d, 0x00, // 0x06google0x3com\0
		0x00, 0x01, // Query Type A (IPv4)
		0x00, 0x01, // Query Class IN (Internet)
	}

	// Send Data
	_, err = conn.Write(data)
	if err != nil {
		t.Fatal("error sending data: %w", err)
	}

	// Receive Response
	buffer := make([]byte, 512) // Standard DNS response size
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal("error receiving response: %w", err)
	}

	t.Log("Response: ", buffer[:n])
}
