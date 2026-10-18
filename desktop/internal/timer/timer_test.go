package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestStopResetsTheTimerAndPreventsRest(t *testing.T) {
	var restStarted atomic.Int32
	tmr := New(100*time.Millisecond, 100*time.Millisecond, Callbacks{
		OnRestStart: func() { restStarted.Add(1) },
	})

	tmr.Start()
	time.Sleep(10 * time.Millisecond)
	tmr.Stop()

	if got := tmr.State(); got != Idle {
		t.Fatalf("state after Stop() = %v, want Idle", got)
	}
	if got, want := tmr.Remaining(), 100*time.Millisecond; got != want {
		t.Fatalf("remaining after Stop() = %v, want %v", got, want)
	}

	// The worker's ticker fires every 50 ms. Waiting past the old work
	// duration proves the cancelled worker cannot advance to a rest cycle.
	time.Sleep(150 * time.Millisecond)
	if got := restStarted.Load(); got != 0 {
		t.Fatalf("rest started %d time(s) after Stop()", got)
	}
}

func TestConfigureUpdatesAnIdleTimer(t *testing.T) {
	tmr := New(time.Minute, 20*time.Second, Callbacks{})
	tmr.Configure(5*time.Minute, 30*time.Second)

	if got, want := tmr.Remaining(), 5*time.Minute; got != want {
		t.Fatalf("remaining after Configure() = %v, want %v", got, want)
	}
}
