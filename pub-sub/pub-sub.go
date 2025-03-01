package pubsub

import (
	"fmt"
	"sync"
)

type Subscription struct {
	Name    string
	Channel chan string
}

func NewSubscription(name string) *Subscription {
	return &Subscription{
		Name:    name,
		Channel: make(chan string, 100),
	}
}

func (s *Subscription) Pull() []string {
	msgs := []string{}

	for {
		select {
		case msg := <-s.Channel:
			msgs = append(msgs, msg)
		default:
			return msgs
		}
	}
}

type PubSub struct {
	subscriptions map[string]map[string]*Subscription
	mu            sync.RWMutex
}

func New() *PubSub {
	return &PubSub{
		subscriptions: make(map[string]map[string]*Subscription),
	}
}

func (ps *PubSub) CreateTopic(topic string) {
	ps.subscriptions[topic] = map[string]*Subscription{}
}

func (ps *PubSub) DeleteTopic(topic string) {
	delete(ps.subscriptions, topic)
}

func (ps *PubSub) Subscribe(topic string, name string) (*Subscription, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := NewSubscription(name)
	existing := ps.subscriptions[topic][name]
	if existing != nil {
		return nil, fmt.Errorf("subscription to topic '%s' with name '%s' already exists", topic, name)
	}
	ps.subscriptions[topic][name] = sub
	return sub, nil
}

func (ps *PubSub) Unsubscribe(topic string, name string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	sub := ps.subscriptions[topic][name]
	if sub == nil {
		return fmt.Errorf("subscription to topic '%s' with name '%s' does not exist", topic, name)
	}
	delete(ps.subscriptions[topic], name)
	return nil
}

func (ps *PubSub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, sub := range ps.subscriptions[topic] {
		select {
		case sub.Channel <- msg:
		default:
			fmt.Printf("Subscription '%s' missed message (too many outstanding messages).\n", sub.Name)
		}
	}
}
