package iote

type Consumer interface {
	// Recv(msg Msg)
	GetID() string
	GetRecvQ() chan Msg
}

type Consumers struct {
	Consumers	map[string][]Consumer
}

var (
	consumers Consumers 
)

func GetConsumers(category string) (cons []Consumer) {
	cons, _ = consumers.Consumers[category]
	return cons
}
