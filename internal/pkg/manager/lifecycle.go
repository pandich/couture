package manager

import (
	"couture/internal/pkg/model"
	"fmt"
)

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (m *busBasedManager) Start() error {
	if m.options.clearScreen {
		clearScreen()
	}
	m.running = true
	for _, poller := range m.pollers {
		m.wg.Add(1)
		go poller(m.wg)
	}
	for _, pusher := range m.pushers {
		if err := pusher.Start(m.wg, func(evt model.Event) {
			m.bus.Publish(eventTopic, pusher, evt)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m *busBasedManager) MustStart() {
	if err := (*m).Start(); err != nil {
		panic(err)
	}
}

func (m *busBasedManager) Stop() {
	m.running = false
	for _, pusher := range m.pushers {
		pusher.Stop()
	}
}

func (m *busBasedManager) Wait() {
	m.wg.Wait()
}
