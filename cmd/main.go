package main

import (
	"log"

	"github.com/brachiGH/firedns/internal/server"
	"github.com/brachiGH/firedns/monitor"
	"github.com/joho/godotenv"
)

func main() {
	// test := flag.String("test", "false", "If the the test flag is set to any value other then false. Testing will begin.")
	// flag.Parse()

	// if *test != "false" {
	// 	fmt.Println("Running tests")
	// 	test.TestDNSOverTLS()
	// 	os.Exit(1)
	// }

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go server.Upd_dns_server()
	go server.ClearCache_Routine()

	var xdp monitor.XDPobj
	err = xdp.Load()
	if err != nil {
		log.Fatal("Failed to load and link XDP: ", err)
	}
	err = xdp.Link()
	if err != nil {
		log.Fatal("Failed to link XDP: ", err)
	}
	defer xdp.UnloadAndCLoseLink()
	go xdp.NICMonitor()
	go func() {
		err := xdp.UpdatePremiumIps()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		err := xdp.UpdateUsageLimitIps()
		if err != nil {
			log.Fatal(err)
		}
	}()

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
