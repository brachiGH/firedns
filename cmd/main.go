package main

import (
	"github.com/brachiGH/firedns/internal/server"
	"github.com/brachiGH/firedns/internal/transport"
	"github.com/brachiGH/firedns/internal/utils/logger"
	"github.com/brachiGH/firedns/monitor"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	log := logger.NewLogger()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go transport.Upd_dns_server()
	go transport.Tcp_dns_server()
	go transport.StartTLSServer("8443")
	go server.ClearCache_Routine()
	go monitor.UpdateQuestions_Routine()

	var xdp monitor.XDPobj
	err = xdp.Load()
	if err != nil {
		log.Fatal("Failed to load and link XDP: ", zap.Error(err))
	}
	err = xdp.Link()
	if err != nil {
		log.Fatal("Failed to link XDP: ", zap.Error(err))
	}
	defer xdp.UnloadAndCLoseLink()
	go xdp.NICMonitor()

	go func() {
		err := xdp.UpdatePremiumIps()
		if err != nil {
			log.Fatal("Updating Premium Ips Failed: ", zap.Error(err))
		}
	}()
	go func() {
		err := xdp.UpdateUsageLimitIps()
		if err != nil {
			log.Fatal("Updating Limited Ips Failed: ", zap.Error(err))
		}
	}()

	select {}
}
