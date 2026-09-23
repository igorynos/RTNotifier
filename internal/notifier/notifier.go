package notifier

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Trigger interface {
	Name() string
	Check(context.Context) (bool, string, error)
}
type Sender interface {
	Send(context.Context, string) error
}
type State interface {
	Ready(context.Context, string, time.Duration) (bool, error)
}

type Runner struct {
	Triggers []Trigger
	Sender   Sender
	State    State
	Cooldown time.Duration
	mu       sync.Mutex
	last     map[string]time.Time
	now      func() time.Time
}

func NewRunner(sender Sender, cooldown time.Duration, triggers ...Trigger) *Runner {
	return &Runner{Triggers: triggers, Sender: sender, Cooldown: cooldown, last: map[string]time.Time{}, now: time.Now}
}

func (r *Runner) RunOnce(ctx context.Context) error {
	for _, trigger := range r.Triggers {
		fired, message, err := trigger.Check(ctx)
		if err != nil {
			return fmt.Errorf("trigger %s: %w", trigger.Name(), err)
		}
		if !fired {
			continue
		}
		if r.State != nil {
			ok, stateErr := r.State.Ready(ctx, trigger.Name(), r.Cooldown)
			if stateErr != nil {
				return fmt.Errorf("state %s: %w", trigger.Name(), stateErr)
			}
			if !ok {
				continue
			}
		} else if !r.ready(trigger.Name()) {
			continue
		}
		if err := r.Sender.Send(ctx, message); err != nil {
			return fmt.Errorf("send %s: %w", trigger.Name(), err)
		}
		r.mark(trigger.Name())
	}
	return nil
}

func (r *Runner) ready(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.now().Sub(r.last[name]) >= r.Cooldown
}
func (r *Runner) mark(name string) { r.mu.Lock(); r.last[name] = r.now(); r.mu.Unlock() }
