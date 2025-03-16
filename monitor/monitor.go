package monitor

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/brachiGH/firedns/monitor/database"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

type XDPobj struct {
	Objs nic_monitorObjects
	link link.Link
}

func (x *XDPobj) LoadAndLink() error {
	// Remove resource limits for kernels <5.11.
	if err := rlimit.RemoveMemlock(); err != nil {
		log.Fatal("Removing memlock:", err)
	}

	// Load the compiled eBPF ELF and load it into the kernel.
	if err := loadNic_monitorObjects(&x.Objs, nil); err != nil {
		log.Fatal("Loading eBPF objects:", err)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Interface on your machine.
	ifname := os.Getenv("ifname")
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Fatalf("Getting interface %s: %s", ifname, err)
	}

	// Attach count_packets to the network interface.
	x.link, err = link.AttachXDP(link.XDPOptions{
		Program:   x.Objs.QueryAnalyser,
		Interface: iface.Index,
	})
	if err != nil {
		log.Fatal("Attaching XDP:", err)
	}

	log.Printf("Monitoring incoming packets on %s..", ifname)
	return nil
}

func (x *XDPobj) UnloadAndCLoseLink() {
	x.link.Close()
	x.Objs.Close()
}

func (x *XDPobj) NICMonitor() {
	//** note: If eBPF/XDP related objects are not defined, execute the 'go generate' command in the directory containing this file. **//

	var db database.Analytics_DB
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	// Periodically fetch the packet counter from PktCount,
	// exit the program when interrupted.
	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-tick:
			var key uint32
			var value uint16
			iter := x.Objs.nic_monitorMaps.QueryCountPerIp.Iterate()
			for iter.Next(&key, &value) {
				fmt.Printf("k: %v v: %+v\n", key, value)

				if value != 0 {
					go func() {
						id, err := db.Update(bson.M{"ip": key}, bson.M{"$inc": bson.M{"QuestionCount": value}})
						fmt.Printf("id: %v err: %s\n", id, err)
					}()
				}

				err := x.Objs.nic_monitorMaps.QueryCountPerIp.Delete(key)
				if err != nil {
					fmt.Println("error deleting:", err)
				}
			}
		case <-stop:
			log.Print("Received signal, exiting..")
			os.Exit(0)
			return
		}
	}
}
