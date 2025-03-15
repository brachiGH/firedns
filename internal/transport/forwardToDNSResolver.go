package transport

import (
	"fmt"
	"net"
)

func ForwardPacketTo(data []byte, serverAddr string) ([]byte, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return []byte{}, fmt.Errorf("error resolving UDP address: %w", err)
	}

	// Create UDP Connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return []byte{}, fmt.Errorf("error connecting to UDP server: %w", err)
	}
	defer conn.Close()

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
