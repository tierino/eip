package pubsub

import (
	"fmt"
	"slices"
	"sync"
	"testing"
)

func TestSubscribe(t *testing.T) {
	ps := New()
	topic := "test_topic"
	msgs := []string{"test message 1", "test message 2", "test message 3"}

	ps.CreateTopic(topic)
	sub, err := ps.Subscribe(topic, "test_subscriber")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	for _, msg := range msgs {
		ps.Publish(topic, msg)
	}

	if got := sub.Pull(); !slices.Equal(got, msgs) {
		t.Fatalf("got %v, want %v", got, msgs)
	}
}

func TestMultipleSubscriptionsSameTopic(t *testing.T) {
	ps := New()
	topic := "test_topic"
	msg := "test message"

	ps.CreateTopic(topic)
	sub1, err := ps.Subscribe(topic, "test_subscriber1")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	sub2, err := ps.Subscribe(topic, "test_subscriber2")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	sub3, err := ps.Subscribe(topic, "test_subscriber3")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	ps.Publish(topic, msg)

	if got := sub1.Pull(); !slices.Equal(got, []string{msg}) {
		t.Fatalf("subcriber 1 got %v, want %v", got, msg)
	}

	if got := sub2.Pull(); !slices.Equal(got, []string{msg}) {
		t.Fatalf("subscriber 2 got %v, want %v", got, msg)
	}

	if got := sub3.Pull(); !slices.Equal(got, []string{msg}) {
		t.Fatalf("subscriber 3 got %v, want %v", got, msg)
	}
}

func TestEnforceSubscriptionsUniquelyNamed(t *testing.T) {
	ps := New()
	topic := "test_topic"

	ps.CreateTopic(topic)
	_, err := ps.Subscribe(topic, "test_subscriber")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}
	_, err = ps.Subscribe(topic, "test_subscriber")
	if err == nil {
		t.Fatalf("expected error creating subscription")
	}
}

func TestMultipleSubscriptionsDifferentTopics(t *testing.T) {
	ps := New()
	topic1 := "test_topic1"
	topic2 := "test_topic2"
	msg := "test message"

	ps.CreateTopic(topic1)
	sub1, err := ps.Subscribe(topic1, "test_subscriber")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	ps.CreateTopic(topic2)
	sub2, err := ps.Subscribe(topic2, "test_subscriber")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	ps.Publish(topic1, msg)
	ps.Publish(topic2, msg)

	if got := sub1.Pull(); !slices.Equal(got, []string{msg}) {
		t.Fatalf("got %v, want %v", got, msg)
	}

	if got := sub2.Pull(); !slices.Equal(got, []string{msg}) {
		t.Fatalf("got %v, want %v", got, msg)
	}
}

func TestUnsubscribe(t *testing.T) {
	ps := New()
	topic := "test_topic"
	msg := "test message"

	ps.CreateTopic(topic)
	sub, err := ps.Subscribe(topic, "test_subscriber")
	if err != nil {
		t.Fatalf("error creating subscription %e", err)
	}

	ps.Unsubscribe(topic, "test_subscriber")
	ps.Publish(topic, msg)

	if got := sub.Pull(); !slices.Equal(got, []string{}) {
		t.Fatalf("got %v, want %v", got, []string{})
	}
}

func TestSubscribeConcurrency(t *testing.T) {
	wantedSubscriptions := 10

	ps := New()
	topic := "test_topic"
	ps.CreateTopic(topic)

	var wg sync.WaitGroup

	wg.Add(wantedSubscriptions)

	for i := range wantedSubscriptions {
		go func() {
			defer wg.Done()
			ps.Subscribe(topic, fmt.Sprintf("test_subscriber%d", i))
		}()
	}

	wg.Wait()

	if got := len(ps.subscriptions[topic]); got != wantedSubscriptions {
		t.Fatalf("got %v, want %v", got, wantedSubscriptions)
	}
}

func TestPublishConcurrency(t *testing.T) {
	wantedMessages := 10

	ps := New()
	topic := "test_topic"
	ps.CreateTopic(topic)

	var wg sync.WaitGroup

	wg.Add(wantedMessages)

	sub, _ := ps.Subscribe(topic, "test_subscriber")

	for i := range wantedMessages {
		go func() {
			defer wg.Done()
			ps.Publish(topic, fmt.Sprintf("test message %d", i))
		}()
	}

	wg.Wait()

	if got := len(sub.Pull()); got != wantedMessages {
		t.Fatalf("got %v, want %v", got, wantedMessages)
	}
}
