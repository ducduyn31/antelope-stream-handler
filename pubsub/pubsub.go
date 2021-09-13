package pubsub

import (
	"sync"
)

func CreatePubSub() *PubSub {
	ps := &PubSub{}
	ps.subs = make(map[string]chan int)
	return ps
}

type PubSub struct {
	mu     sync.RWMutex
	subs   map[string]chan int
	closed bool
}

func (ps *PubSub) Subscribe(source string) <-chan int{
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan int, 1)
	ps.subs[source] = ch
	return ch
}

func (ps *PubSub) Signal(source string, value int) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}
	ps.subs[source] <- value
}

func (ps *PubSub) CloseAndRemoveAll()  {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for source, sub := range ps.subs {
			close(sub)
			delete(ps.subs, source)
		}
	}
}

func (ps *PubSub) CloseOne(source string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	close(ps.subs[source])
}

func (ps *PubSub) RemoveOne(source string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	delete(ps.subs, source)
}
