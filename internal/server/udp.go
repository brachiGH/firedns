package server

import (
	"net"

	"github.com/brachiGH/firedns/internal/utils/logger"
	"go.uber.org/zap"
)

func Upd_dns_server() {
	log := logger.NewLogger()
	defer func() {
		if err := log.Sync(); err != nil {
			panic(err)
		}
	}()

	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:2053")
	if err != nil {
		log.Fatal("Failed to resolve UDP address: ", zap.Error(err))
		panic(err)
	}
	log.Info("Server started", zap.String("port", "2053"))

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("Failed to bind to address: ", zap.Error(err))
		panic(err)
	}

	defer func() {
		if err := udpConn.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, maxPacketsize+1) // max size is 512, The extra byte is to detect if the packet is lager then 512bytes
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal("Error receiving data: ", zap.Error(err))
		}
		if size == maxPacketsize+1 {
			// UDP response larger than 512 bytes are not supported
			continue
		}

		data, err := handle(buf[:size])
		if err != nil {
			log.Fatal("Fail to handle request: ", zap.Error(err))
		}

		_, err = udpConn.WriteToUDP(data, source)
		if err != nil {
			log.Fatal("Failed to send response: ", zap.Error(err))
		}
	}
}
