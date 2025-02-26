package dynamicrouter

type DynamicRouter struct {
	UnknownQueue *MessageQueue
	RouteTable   map[string]*MessageQueue
}

func New() *DynamicRouter {
	return &DynamicRouter{
		RouteTable: make(map[string]*MessageQueue),
		UnknownQueue: &MessageQueue{
			Name:  "/unknown",
			queue: make(chan string, 10),
		},
	}
}

func (r *DynamicRouter) HandleMessage(msg *Message) {
	if _, ok := r.RouteTable[msg.Route]; ok {
		r.RouteTable[msg.Route].Enqueue(msg.Content)
	} else {
		r.UnknownQueue.Enqueue(msg.Content)
	}
}

func (r *DynamicRouter) HandleControlMessage(msg *ControlMessage) {
	if msg.Action == RemoveRoute {
		r.removeRoute(msg.Route)
		return
	}
	if _, ok := r.RouteTable[msg.Route]; ok {
		return
	}
	r.addRoute(msg.Route, msg.QueueName)
}

func (r *DynamicRouter) addRoute(route string, name string) {
	queue := r.findQueue(name)
	if queue != nil {
		r.RouteTable[route] = queue
		return
	}

	r.RouteTable[route] = &MessageQueue{
		Name:  name,
		queue: make(chan string, 10),
	}
}

func (r *DynamicRouter) removeRoute(route string) {
	delete(r.RouteTable, route)
}

func (r *DynamicRouter) findQueue(name string) *MessageQueue {
	for _, queue := range r.RouteTable {
		if queue.Name == name {
			return queue
		}
	}
	return nil
}
