package monitor

import (
	"log"
	"time"

	"github.com/brachiGH/firedns/internal/utils/config"
	"github.com/brachiGH/firedns/monitor/database"
)

func (x *XDPobj) UpdateUsageLimitIps() error {
	var db database.Analytics_DB
	err := db.Connect()
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Disconnect(); err != nil {
			log.Fatalln(err)
		}
	}()

	tick := time.Tick(config.UpdateUsageLimitIps__TickDuration)
	for range tick {
		keys := []uint32{}
		values := []uint32{}
		_, err := x.Objs.UsageLimitExededIps.BatchUpdate(keys, values, nil)
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
		return err
	}
	defer func() {
		if err := db.Disconnect(); err != nil {
			log.Fatalln(err)
		}
	}()

	tick := time.Tick(config.UpdatePremiumIps__TickDuration)
	for range tick {
		keys := []uint32{}
		values := []uint32{}
		_, err := x.Objs.PremiumIps.BatchUpdate(keys, values, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
