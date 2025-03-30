package transport

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/brachiGH/firedns/internal/server"
	"github.com/brachiGH/firedns/internal/utils/config"
	"github.com/brachiGH/firedns/internal/utils/logger"
	"go.uber.org/zap"
)

func Tcp_dns_server() {
	log := logger.NewLogger()

	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:2053")
	if err != nil {
		log.Fatal("Failed to resolve TCP address: ", zap.Error(err))
	}
	log.Info("TCP Server started", zap.String("port", "2053"))

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("Failed to bind to address: ", zap.Error(err))
	}

	defer func() {
		if err := tcpListener.Close(); err != nil {
			log.Fatal("failed to close listener", zap.Error(err))
		}
	}()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Error("Error accepting connection: ", zap.Error(err))
			continue
		}

		go handleTCPConnection(conn, log)
	}
}

func handleTCPConnection(conn net.Conn, log *zap.Logger) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Error("Error closing connection: ", zap.Error(cerr))
		}
	}()

	// Read message length (first 2 bytes in TCP DNS)
	lengthBuf := make([]byte, 2)
	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		log.Error("Error reading message length: ", zap.Error(err))
		return
	}

	length := binary.BigEndian.Uint16(lengthBuf)
	if length > uint16(server.MaxPacketsize) {
		log.Error("Message too large", zap.Uint16("length", length))
		return
	}

	// Read the actual message
	buf := make([]byte, length)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		log.Error("Error reading message: ", zap.Error(err))
		return
	}

	addr := conn.RemoteAddr().(*net.TCPAddr)
	log.Debug("new connection", zap.String("IP", addr.IP.String()))

	data, err := server.HandleDnsMessage(buf, addr.IP, config.NO_PROFILE_ID)
	if err != nil {
		log.Error("Fail to handle request: ", zap.Error(err))
		return
	}

	// Write response length
	responseLengthBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(responseLengthBuf, uint16(len(data)))
	_, err = conn.Write(responseLengthBuf)
	if err != nil {
		log.Error("Failed to send response length: ", zap.Error(err))
		return
	}

	// Write response
	_, err = conn.Write(data)
	if err != nil {
		log.Error("Failed to send response: ", zap.Error(err))
	}
}
