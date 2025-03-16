package monitor

import (
	"log"
	"time"

	"github.com/brachiGH/firedns/monitor/database"
)

func (x *XDPobj) UpdateUsageLimitIps() error {
	var db database.Analytics_DB
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	tick := time.Tick(time.Minute * 10)
	for _ = range tick {
		keys := []uint32{}
		values := []uint32{}
		_, err := x.Objs.nic_monitorMaps.UsageLimitExededIps.BatchUpdate(keys, values, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (x *XDPobj) UpdatePremiumIps() error {
	var db database.Analytics_DB
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Disconnect()

	tick := time.Tick(time.Minute * 10)
	for _ = range tick {
		keys := []uint32{}
		values := []uint32{}
		_, err := x.Objs.nic_monitorMaps.PremiumIps.BatchUpdate(keys, values, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
