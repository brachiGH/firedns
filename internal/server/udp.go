package server

import (
	"net"

	"github.com/brachiGH/firedns/internal/utils/logger"
	"github.com/brachiGH/firedns/monitor"
	"go.uber.org/zap"
)

func Upd_dns_server() {
	log := logger.NewLogger()

	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:2053")
	if err != nil {
		log.Fatal("Failed to resolve UDP address: ", zap.Error(err))
	}
	log.Info("Server started", zap.String("port", "2053"))

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to bind to address: ", zap.Error(err))
	}

	defer func() {
		if err := udpConn.Close(); err != nil {
			log.Fatal("failed to close connection", zap.Error(err))
		}
	}()

	go monitor.UpdateQuestions_Routine()

	buf := make([]byte, maxPacketsize+1) // max size is 512, The extra byte is to detect if the packet is lager then 512bytes
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Error("Error receiving data: ", zap.Error(err))
			continue
		}
		if size >= maxPacketsize+1 {
			// UDP response larger than 512 bytes are not supported
			continue
		}

		data, err := handle(buf[:size], source.IP)
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
