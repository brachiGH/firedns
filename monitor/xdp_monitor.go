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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type XDPobj struct {
	Objs     nic_monitorObjects
	link     link.Link
	IsLinked *bool
	IsLoaded *bool
}

func (x *XDPobj) Load() error {
	if x.IsLoaded == nil || !*x.IsLoaded {
		// Remove resource limits for kernels <5.11.
		if err := rlimit.RemoveMemlock(); err != nil {
			return fmt.Errorf("removing memlock: %w", err)
		}

		// Load the compiled eBPF ELF and load it into the kernel.
		if err := loadNic_monitorObjects(&x.Objs, nil); err != nil {
			return fmt.Errorf("loading eBPF objects: %w", err)
		}

		var loaded = true
		x.IsLoaded = &loaded

		return nil
	}

	return fmt.Errorf("XDP object is already loaded")
}

func (x *XDPobj) Link() error {
	if x.IsLinked == nil || !*x.IsLinked {
		// Interface on your machine.
		ifname := os.Getenv("ifname")
		iface, err := net.InterfaceByName(ifname)
		if err != nil {
			return fmt.Errorf("getting interface %s: %w", ifname, err)
		}

		// Attach count_packets to the network interface.
		x.link, err = link.AttachXDP(link.XDPOptions{
			Program:   x.Objs.QueryAnalyser,
			Interface: iface.Index,
		})

		if err != nil {
			return fmt.Errorf("attaching XDP: %w", err)
		}

		log.Printf("Monitoring incoming packets on %s..", ifname)
		var linked = true
		x.IsLinked = &linked

		return nil
	}

	return fmt.Errorf("XDP object is already linked")
}

func (x *XDPobj) UnloadAndCLoseLink() {
	err := x.link.Close()
	if err != nil {
		log.Printf("Error closing link: %v", err)
	}

	err = x.Objs.Close()
	if err != nil {
		log.Printf("Error closing objects: %v", err)
	}
}

func (x *XDPobj) NICMonitor() {
	//** note: If eBPF/XDP related objects are not defined, execute the 'go generate' command in the directory containing this file. **//

	var db database.Analytics_DB
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Disconnect(); err != nil {
			panic(err)
		}
	}()

	// Periodically fetch the packet counter from PktCount,
	// exit the program when interrupted.
	tick := time.Tick(time.Second)
	stop := make(chan os.Signal, 5)
	signal.Notify(stop, os.Interrupt)
	for {
		select {
		case <-tick:
			var key uint32
			var value uint32
			iter := x.Objs.QueryCountPerIp.Iterate()
			for iter.Next(&key, &value) {
				log.Printf("k: %v v: %+v\n", key, value)
				updates := []mongo.WriteModel{}

				if value != 0 {
					updates = append(updates, mongo.NewUpdateOneModel().
						SetFilter(bson.M{"ip": key}).
						SetUpdate(bson.M{"$inc": bson.M{"QuestionCount": value}}).SetUpsert(true))
				}
				go func() {
					err := db.UpdateMany(updates)
					if err != nil {
						log.Println("error updating:", err)
					}
				}()

				err := x.Objs.QueryCountPerIp.Delete(key)
				if err != nil {
					log.Println("error deleting:", err)
				}
			}
			if err := iter.Err(); err != nil {
				log.Printf("Iteration error: %v", err)
			}
		case <-stop:
			log.Print("Received signal, exiting..")
			os.Exit(0)
			return
		}
	}
}
