package dynamicrouter

import (
	"testing"
)

func TestHandleMessageWithoutRoute(t *testing.T) {
	router := New()

	content := "test message"
	msg := &Message{Route: "route1", Content: content}

	router.HandleMessage(msg)

	if got := router.UnknownQueue.Dequeue(); got != content {
		t.Fatalf("got %v, want %v", got, content)
	}
}

func TestHandleMessageWithRoute(t *testing.T) {
	router := New()

	queueName := "/queue1"
	ctrlMsg := &ControlMessage{Action: AddRoute, Route: "route1", QueueName: queueName}
	content := "test message"
	msg := &Message{Route: "route1", Content: content}

	router.HandleControlMessage(ctrlMsg)
	router.HandleMessage(msg)

	if got := router.findQueue(queueName).Dequeue(); got != content {
		t.Fatalf("got %v, want %v", got, content)
	}
}

func TestControlMessageRemoveRoute(t *testing.T) {
	router := New()

	queueName := "/queue1"
	addCtrlMsg := &ControlMessage{Action: AddRoute, Route: "route1", QueueName: queueName}
	router.HandleControlMessage(addCtrlMsg)

	removeCtrlMsg := &ControlMessage{Action: RemoveRoute, Route: "route1", QueueName: queueName}
	router.HandleControlMessage(removeCtrlMsg)

	if got := router.findQueue(queueName); got != nil {
		t.Fatalf("got %v, want %v", got, nil)
	}
}
