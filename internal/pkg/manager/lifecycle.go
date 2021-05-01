package manager

import (
	"fmt"
)

const cls = "\033[H\033[2J"

func (m *busBasedManager) Start() error {
	if m.options.clearScreen {
		fmt.Print(cls)
	}
	m.running = true
	for _, poller := range m.pollers {
		m.wg.Add(1)
		go poller(m.wg)
	}
	for _, pusher := range m.pushers {
		m.wg.Add(1)
		if err := pusher.Start(m.wg); err != nil {
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
