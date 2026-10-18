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
	runID        uint64
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
	t.runID++
	runID := t.runID
	stopCh := t.stopCh
	t.remaining = t.workDuration
	t.state = Running
	onStateChange := t.cb.OnStateChange
	t.mu.Unlock()

	if onStateChange != nil {
		onStateChange(Running)
	}
	go t.run(runID, stopCh)
}

func (t *Timer) Stop() {
	t.mu.Lock()
	if t.state == Idle {
		t.mu.Unlock()
		return
	}
	close(t.stopCh)
	t.stopCh = nil
	t.runID++ // Invalidates an in-flight timer goroutine.
	t.state = Idle
	t.remaining = t.workDuration
	remaining := t.remaining
	onStateChange := t.cb.OnStateChange
	onTick := t.cb.OnTick
	t.mu.Unlock()

	if onStateChange != nil {
		onStateChange(Idle)
	}
	if onTick != nil {
		onTick(remaining, Idle)
	}
}

func (t *Timer) run(runID uint64, stopCh <-chan struct{}) {
	for {
		// Work phase
		if !t.countdown(t.durationFor(Running, runID), Running, runID, stopCh) {
			return
		}

		// Rest phase
		if !t.setState(Resting, runID, stopCh) {
			return
		}
		if t.isActive(runID, stopCh) && t.cb.OnRestStart != nil {
			t.cb.OnRestStart()
		}

		if !t.countdown(t.durationFor(Resting, runID), Resting, runID, stopCh) {
			return
		}

		if t.isActive(runID, stopCh) && t.cb.OnRestEnd != nil {
			t.cb.OnRestEnd()
		}

		// Next cycle
		if !t.setState(Running, runID, stopCh) {
			return
		}
	}
}

func (t *Timer) durationFor(phase State, runID uint64) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.runID != runID {
		return 0
	}
	if phase == Resting {
		return t.restDuration
	}
	return t.workDuration
}

func (t *Timer) countdown(duration time.Duration, phase State, runID uint64, stopCh <-chan struct{}) bool {
	if duration <= 0 {
		return false
	}
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
			if t.runID != runID || t.stopCh != stopCh {
				t.mu.Unlock()
				return false
			}
			t.remaining = rem
			onTick := t.cb.OnTick
			t.mu.Unlock()
			if onTick != nil {
				onTick(rem, phase)
			}
			if rem <= 0 {
				return true
			}
		}
	}
}

func (t *Timer) setState(state State, runID uint64, stopCh <-chan struct{}) bool {
	t.mu.Lock()
	if t.runID != runID || t.stopCh != stopCh {
		t.mu.Unlock()
		return false
	}
	t.state = state
	onStateChange := t.cb.OnStateChange
	t.mu.Unlock()

	if onStateChange != nil {
		onStateChange(state)
	}
	return true
}

func (t *Timer) isActive(runID uint64, stopCh <-chan struct{}) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.runID == runID && t.stopCh == stopCh
}
