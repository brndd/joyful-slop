package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/stretchr/testify/require"
)

type timedRuleStub struct {
	calls int
}

func (rule *timedRuleStub) TimerEvents(_ *string) []mappingrules.OutputEvent {
	rule.calls++
	return nil
}

func TestTimerWatcherQueuesRuleWithoutEvaluatingIt(t *testing.T) {
	rule := &timedRuleStub{}
	channel := make(chan ChannelEvent, 1)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go timerWatcher(rule, channel, ctx, &wg)

	select {
	case event := <-channel:
		require.Same(t, rule, event.Rule)
	case <-time.After(time.Second):
		t.Fatal("timer watcher did not queue a rule tick")
	}
	require.Zero(t, rule.calls)
	cancel()
	wg.Wait()
}
