package poller

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestPoll_ImmediateSuccess(t *testing.T) {
	cfg := PollConfig{InitialInterval: 10 * time.Millisecond, MaxInterval: 100 * time.Millisecond, Timeout: 5 * time.Second, BackoffFactor: 1.5}
	calls := 0
	data, err := Poll(context.Background(), cfg, func() PollResult {
		calls++
		return PollResult{Done: true, Data: "ok"}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != "ok" {
		t.Errorf("expected 'ok', got %v", data)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestPoll_EventualSuccess(t *testing.T) {
	cfg := PollConfig{InitialInterval: 10 * time.Millisecond, MaxInterval: 50 * time.Millisecond, Timeout: 5 * time.Second, BackoffFactor: 1.5}
	calls := 0
	data, err := Poll(context.Background(), cfg, func() PollResult {
		calls++
		if calls >= 3 {
			return PollResult{Done: true, Data: "done"}
		}
		return PollResult{Done: false}
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data != "done" {
		t.Errorf("expected 'done', got %v", data)
	}
	if calls < 3 {
		t.Errorf("expected at least 3 calls, got %d", calls)
	}
}

func TestPoll_Timeout(t *testing.T) {
	cfg := PollConfig{InitialInterval: 50 * time.Millisecond, MaxInterval: 50 * time.Millisecond, Timeout: 200 * time.Millisecond, BackoffFactor: 1.0}
	_, err := Poll(context.Background(), cfg, func() PollResult {
		return PollResult{Done: false}
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestPoll_ErrorPropagation(t *testing.T) {
	cfg := PollConfig{InitialInterval: 10 * time.Millisecond, MaxInterval: 100 * time.Millisecond, Timeout: 5 * time.Second, BackoffFactor: 1.5}
	calls := 0
	_, err := Poll(context.Background(), cfg, func() PollResult {
		calls++
		if calls >= 2 {
			return PollResult{Err: fmt.Errorf("task failed")}
		}
		return PollResult{Done: false}
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "task failed" {
		t.Errorf("expected 'task failed', got '%s'", err.Error())
	}
}
