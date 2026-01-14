package notifier

import (
	"sync"

	"github.com/google/uuid"
)

type subscription struct {
	ch chan struct{}
	wg sync.WaitGroup
}

type Notifier struct {
	mu            sync.Mutex
	subscriptions map[string]*subscription
}

func (n *Notifier) Subscribe() (<-chan struct{}, func()) {
	id := uuid.NewString()
	ch := make(chan struct{})

	n.mu.Lock()
	defer n.mu.Unlock()

	if n.subscriptions == nil {
		n.subscriptions = make(map[string]*subscription)
	}
	n.subscriptions[id] = &subscription{ch: ch}

	return ch, func() {
		n.mu.Lock()
		defer n.mu.Unlock()

		sub, ok := n.subscriptions[id]
		if ok {
			// wait for any in-flight notifications to finish
			sub.wg.Wait()
			close(sub.ch)
			delete(n.subscriptions, id)
		}
	}
}

func (n *Notifier) Notify() {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, sub := range n.subscriptions {
		sub.wg.Go(func() {
			sub.ch <- struct{}{}
		})
	}
}
