package dynamicrouter

type ControlMessageAction string

const (
	AddRoute    ControlMessageAction = "add"
	RemoveRoute ControlMessageAction = "remove"
)

type ControlMessage struct {
	Action    ControlMessageAction
	Route     string
	QueueName string
}

type Message struct {
	Route   string
	Content string
}
