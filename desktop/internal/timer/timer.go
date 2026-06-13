package timer

import (
	"sync"
	"time"
)

type State int

const (
	Idle State = iota
	Running
	Resting
)

type Callbacks struct {
	OnTick        func(remaining time.Duration, state State)
	OnRestStart   func()
	OnRestEnd     func()
	OnStateChange func(State)
}

type Timer struct {
	mu           sync.Mutex
	state        State
	workDuration time.Duration
	restDuration time.Duration
	remaining    time.Duration
	stopCh       chan struct{}
	cb           Callbacks
}

func New(work, rest time.Duration, cb Callbacks) *Timer {
	return &Timer{
		state:        Idle,
		workDuration: work,
		restDuration: rest,
		remaining:    work,
		cb:           cb,
	}
}

func (t *Timer) State() State {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state
}

func (t *Timer) Remaining() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.remaining
}

func (t *Timer) Configure(work, rest time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.workDuration = work
	t.restDuration = rest
	if t.state == Idle {
		t.remaining = work
	}
}

func (t *Timer) Start() {
	t.mu.Lock()
	if t.state != Idle {
		t.mu.Unlock()
		return
	}
	t.stopCh = make(chan struct{})
	t.remaining = t.workDuration
	t.setState(Running)
	t.mu.Unlock()

	go t.run()
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == Idle {
		return
	}
	close(t.stopCh)
	t.setState(Idle)
	t.remaining = t.workDuration
	if t.cb.OnTick != nil {
		t.cb.OnTick(t.remaining, Idle)
	}
}

func (t *Timer) setState(s State) {
	t.state = s
	if t.cb.OnStateChange != nil {
		t.cb.OnStateChange(s)
	}
}

func (t *Timer) run() {
	for {
		// Work phase
		if !t.countdown(t.workDuration, Running) {
			return
		}

		// Rest phase
		t.mu.Lock()
		t.setState(Resting)
		t.mu.Unlock()
		if t.cb.OnRestStart != nil {
			t.cb.OnRestStart()
		}

		if !t.countdown(t.restDuration, Resting) {
			return
		}

		if t.cb.OnRestEnd != nil {
			t.cb.OnRestEnd()
		}

		// Next cycle
		t.mu.Lock()
		t.setState(Running)
		t.mu.Unlock()
	}
}

func (t *Timer) countdown(duration time.Duration, phase State) bool {
	start := time.Now()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopCh:
			return false
		case <-ticker.C:
			elapsed := time.Since(start)
			rem := duration - elapsed
			if rem < 0 {
				rem = 0
			}
			t.mu.Lock()
			t.remaining = rem
			t.mu.Unlock()
			if t.cb.OnTick != nil {
				t.cb.OnTick(rem, phase)
			}
			if rem <= 0 {
				return true
			}
		}
	}
}
