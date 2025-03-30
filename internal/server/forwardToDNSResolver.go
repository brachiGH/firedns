package server

import (
	"fmt"
	"log"
	"net"

	"github.com/brachiGH/firedns/internal/utils/config"
)

func ForwardPacketTo(data []byte) ([]byte, error) {
	addr, err := net.ResolveUDPAddr("udp", config.UDP_ns_addr)
	if err != nil {
		return []byte{}, fmt.Errorf("error resolving UDP address: %w", err)
	}

	// Create UDP Connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return []byte{}, fmt.Errorf("error connecting to UDP server: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close name server connection: %v", err)
		}
	}()

	// Send Data
	_, err = conn.Write(data)
	if err != nil {
		return []byte{}, fmt.Errorf("error sending data: %w", err)
	}

	// Receive Response
	buffer := make([]byte, 512) // Standard DNS response size
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return []byte{}, fmt.Errorf("error receiving response: %w", err)
	}

	return buffer[:n], nil
}
