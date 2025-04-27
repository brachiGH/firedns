package monitor

import (
	"time"

	"github.com/brachiGH/firedns/internal/utils/config"
)

func (x *XDPobj) UpdateUsageLimitIps() error {
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
