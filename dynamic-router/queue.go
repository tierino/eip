package dynamicrouter

type MessageQueue struct {
	Name  string
	queue chan string
}

func (q *MessageQueue) Enqueue(msg string) {
	q.queue <- msg
}

func (q *MessageQueue) Dequeue() string {
	return <-q.queue
}
