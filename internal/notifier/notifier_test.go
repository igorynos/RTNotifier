package notifier

import (
	"context"
	"testing"
	"time"
)

type trigger struct{}

func (trigger) Name() string                                { return "usage" }
func (trigger) Check(context.Context) (bool, string, error) { return true, "limit reached", nil }

type sender struct{ calls int }

func (s *sender) Send(context.Context, string) error { s.calls++; return nil }
func TestCooldownPreventsDuplicates(t *testing.T) {
	s := &sender{}
	r := NewRunner(s, time.Hour, trigger{})
	now := time.Unix(10000, 0)
	r.now = func() time.Time { return now }
	if err := r.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := r.RunOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s.calls != 1 {
		t.Fatalf("got %d sends", s.calls)
	}
}
