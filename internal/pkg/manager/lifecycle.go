package manager

import (
	"couture/internal/pkg/model"
)

func (mgr *busBasedManager) Start() error {
	mgr.running = true
	for _, poller := range mgr.pollers {
		mgr.wg.Add(1)
		go poller(mgr.wg)
	}
	for _, pusher := range mgr.pushers {
		if err := pusher.Start(mgr.wg, func(event model.Event) {
			mgr.bus.Publish(eventTopic, pusher, event)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *busBasedManager) MustStart() {
	if err := (*mgr).Start(); err != nil {
		panic(err)
	}
}

func (mgr *busBasedManager) Stop() {
	mgr.running = false
	for _, pusher := range mgr.pushers {
		pusher.Stop()
	}
}

func (mgr *busBasedManager) Wait() {
	mgr.wg.Wait()
}
