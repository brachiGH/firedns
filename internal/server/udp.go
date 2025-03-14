package server

import (
	"net"

	"github.com/brachiGH/firedns/internal/utils"
)

func Upd_dns_server() {
	log := utils.NewLogger()

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		log.Error.Fatalln("Failed to resolve UDP address: ", err)
		return
	}
	log.Info.Println("Server started", "port", 53)

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error.Fatalln("Failed to bind to address: ", err)
		return
	}
	defer udpConn.Close()

	var resolver *net.Resolver = newResolver(ns)

	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			log.Error.Fatalln("Error receiving data: ", err)
			return
		}
		msg, err := NewDNSMessage(buf[:size])
		if err != nil {
			log.Error.Fatalln("Error parsing incoming message: ", err)
			return
		}
		response, err := handle(resolver, msg)
		if err != nil {
			log.Error.Fatalln("Fail to handle request: ", err)
			return
		}
		_, err = udpConn.WriteToUDP(response.AsBytes(), source)
		if err != nil {
			log.Error.Fatalln("Failed to send response: ", err)
			return
		}
	}
}
