package notifications

import (
	"sync"
)

// Broker is a struct holding information about listening subscriptions
type Broker struct {
	subscribers subscribers
	channels    map[string]subscribers
	mut         sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		channels: map[string]subscribers{},
	}
}

// Publish adds a new notification to the appropriate channel based on topic
func (b *Broker) Publish(topic, msg string) {
	b.mut.RLock()
	channels := b.channels[topic]
	b.mut.RUnlock()
	for _, s := range channels {
		go s.Push(newNotification(msg, topic))
	}
}

func (b *Broker) Subscribe(topic string) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if b.channels[topic] == nil {
		b.channels[topic] = subscribers{}
	}
	// TODO: creates topic but doesn't actually have a subscriber implementation
}
