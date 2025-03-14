package main

import "github.com/brachiGH/firedns/internal/server"

func main() {
	go server.Upd_dns_server()

	select {}
}
