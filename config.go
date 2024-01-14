package otto

type Config interface {
	GetAddr() string
	GetAppdir() string
	GetBroker() string
}
