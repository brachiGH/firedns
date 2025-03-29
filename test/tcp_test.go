package test

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/brachiGH/firedns/internal/transport"
	"github.com/brachiGH/firedns/internal/utils/config"
)

func TestDNSTCPMessage(t *testing.T) {
	go transport.Tcp_dns_server()

	// Wait for the TCP server to start
	time.Sleep(1 * time.Second)

	// Connect to TCP server
	conn, err := net.Dial("tcp", config.TCP_local_ns_addr)
	if err != nil {
		t.Fatal("error connecting to TCP server:", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			t.Logf("Failed to close connection: %v", cerr)
		}
	}()

	// Minimal DNS Query (Query for google.com, type A)
	dnsQuery := []byte{
		0x00, 0x01, // Transaction ID
		0x00, 0x00, // Flags
		0x00, 0x01, // Question Count
		0x00, 0x00, // Answer RRs
		0x00, 0x00, // Authority RRs
		0x00, 0x00, // Additional RRs
		0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, // 0x06google0x3com\0
		0x00, 0x01, // Query Type A (IPv4)
		0x00, 0x01, // Query Class IN (Internet)
	}

	// Create TCP message (length prefix + DNS query)
	tcpMessage := make([]byte, 2+len(dnsQuery))
	binary.BigEndian.PutUint16(tcpMessage[0:2], uint16(len(dnsQuery)))
	copy(tcpMessage[2:], dnsQuery)

	// Send Data
	_, err = conn.Write(tcpMessage)
	if err != nil {
		t.Fatal("error sending data:", err)
	}

	// Read response length (2 bytes)
	lengthBuf := make([]byte, 2)
	_, err = conn.Read(lengthBuf)
	if err != nil {
		t.Fatal("error reading response length:", err)
	}
	responseLen := binary.BigEndian.Uint16(lengthBuf)

	// Read response
	response := make([]byte, responseLen)
	_, err = conn.Read(response)
	if err != nil {
		t.Fatal("error reading response:", err)
	}

	t.Logf("Response length: %d bytes", responseLen)
	t.Logf("Response: %x", response)
}
