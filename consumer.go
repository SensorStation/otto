package main

type Consumer interface {
	// Recv(msg Msg)
	GetID() string
	GetRecvQ() chan Msg
}

