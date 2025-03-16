package main

import (
	"github.com/brachiGH/firedns/internal/server"
	"github.com/brachiGH/firedns/monitor"
)

func main() {
	// test := flag.String("test", "false", "If the the test flag is set to any value other then false. Testing will begin.")
	// flag.Parse()

	// if *test != "false" {
	// 	fmt.Println("Running tests")
	// 	test.TestDNSOverTLS()
	// 	os.Exit(1)
	// }

	go server.Upd_dns_server()

	var xdp monitor.XDPobj
	xdp.LoadAndLink()
	defer xdp.UnloadAndCLoseLink()
	go xdp.NICMonitor()
	go xdp.UpdatePremiumIps()
	go xdp.UpdateUsageLimitIps()
	// go transport.TLS()
	// go server.StartDoTServer("/etc/letsencrypt/live/brachi.me/fullchain.pem", "/etc/letsencrypt/live/brachi.me/privkey.pem")

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Log the directory being served
	// fmt.Println("Serving files from:", dir)

	// // Serve files in the current directory (or subdirectories)
	// http.Handle("/", http.FileServer(http.Dir(dir)))

	// // Start the server on port 8080
	// fmt.Println("Server is running on http://localhost:8080")
	// err = http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal("Server failed to start: ", err)
	// }

	select {}
}
