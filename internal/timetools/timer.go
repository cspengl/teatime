package timetools

import "time"

type Timer struct {
	C                chan time.Duration
	running          bool
	stop             chan struct{}
	ticker           *time.Ticker
	initial, current time.Duration
}

func NewTimer(duration time.Duration) *Timer {
	return &Timer{
		C:       make(chan time.Duration),
		stop:    make(chan struct{}, 1),
		running: false,
		current: duration,
		initial: duration,
	}
}

func (t *Timer) Start() chan time.Duration {
	t.ticker = time.NewTicker(1 * time.Second)
	if t.current > 0 {
		go func() {
			for {
				select {
				case <-t.ticker.C:
					t.current -= 1 * time.Second
					t.C <- t.current
					if t.current == 0 {
						t.Stop()
					}
				case <-t.stop:
					t.running = false
					return
				}
			}
		}()
		t.running = true
	}
	return t.C
}

func (t *Timer) Stop() {
	t.stop <- struct{}{}
}

func (t *Timer) Get() time.Duration {
	return t.current
}

func (t *Timer) Set(d time.Duration) {
	t.current = d
}

func (t *Timer) IsRunning() bool {
	return t.running
}

func (t *Timer) Reset() {
	t.current = t.initial
}
