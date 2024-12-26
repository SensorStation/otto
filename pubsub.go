package otto

import "github.com/sensorstation/otto/message"

// Subscriber is an interface that defines a struct needs to have the
// Callback(topic string, data []byte) function defined.
type Subscriber interface {
	Callback(msg *message.Msg)
}

type Subscribed func(msg *message.Msg)

// Publisher interface allows objects to publish message to a particular
// topic as defined in the message.Msg
type Publisher interface {
	Publish(msg *message.Msg)
}
