package messanger

// Subscriber is an interface that defines a struct needs to have the
// Callback(topic string, data []byte) function defined.
type MsgHandle func(msg *Msg)

// Subscriber interface defines a type that can subscribe to published messages
type Subscriber interface {
	Callback(msg *Msg)
}

// Publisher interface allows objects to publish message to a particular
// topic as defined in the message.Msg
type Publisher interface {
	Publish(msg *Msg)
}

// Messanger represents a type that can publish and subscribe to messages
type Messanger interface {
	Publisher()
	Subscriber()
}
