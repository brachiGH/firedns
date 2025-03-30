package transport

import (
	"net"

	"github.com/brachiGH/firedns/internal/server"
	"github.com/brachiGH/firedns/internal/utils/config"
	"github.com/brachiGH/firedns/internal/utils/logger"
	"go.uber.org/zap"
)

func Upd_dns_server() {
	log := logger.NewLogger()

	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:2053")
	if err != nil {
		log.Fatal("Failed to resolve UDP address: ", zap.Error(err))
	}
	log.Info("UDP Server started", zap.String("port", "2053"))

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to bind to address: ", zap.Error(err))
	}

	defer func() {
		if err := udpConn.Close(); err != nil {
			log.Fatal("failed to close connection", zap.Error(err))
		}
	}()

	buf := make([]byte, server.MaxPacketsize+1) // max size is 512, The extra byte is to detect if the packet is lager then 512bytes
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Error("Error receiving data: ", zap.Error(err))
			continue
		}
		log.Debug("new connection", zap.Any("IP", source.IP.String()))
		if size >= server.MaxPacketsize+1 {
			// UDP response larger than 512 bytes are not supported
			continue
		}

		data, err := server.HandleDnsMessage(buf[:size], source.IP, config.NO_PROFILE_ID)
		if err != nil {
			log.Error("Fail to handle request: ", zap.Error(err))
			continue
		}

		_, err = udpConn.WriteToUDP(data, source)
		if err != nil {
			log.Error("Failed to send response: ", zap.Error(err))
		}
	}
}
