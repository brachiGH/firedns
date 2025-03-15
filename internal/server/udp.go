package server

import (
	"net"

	"github.com/brachiGH/firedns/internal/utils"
)

func Upd_dns_server() {
	log := utils.NewLogger()

	udpAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:2053")
	if err != nil {
		log.Error.Fatalln("Failed to resolve UDP address: ", err)
		return
	}
	log.Info.Println("Server started", "port", 2053)

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error.Fatalln("Failed to bind to address: ", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, maxPacketsize+1) // max size is 512, The extra byte is to detect if the packet is lager then 512bytes
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Error.Fatalln("Error receiving data: ", err)
			return
		}
		if size == maxPacketsize+1 {
			// UDP response larger than 512 bytes are not supported
			continue
		}

		qs, err := NewDNSMessage(buf[:size])
		if err != nil {
			log.Error.Fatalln("Error parsing incoming message: ", err)
			return
		}

		data, err := handle(buf[:size], qs)
		if err != nil {
			log.Error.Fatalln("Fail to handle request: ", err)
			return
		}
		_, err = udpConn.WriteToUDP(data, source)
		if err != nil {
			log.Error.Fatalln("Failed to send response: ", err)
			return
		}
	}
}
