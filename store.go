package iote

import (
	"fmt"
	"log"
)

type Store interface {
	Store(msg *Msg) error
}

type MsgStore struct {
	Source map[string]map[string]float64
}

var (
	store *MsgStore
)

func init() {
	store = NewStore()
}

func NewStore() *MsgStore {
	m := make(map[string]map[string]float64)
	store = &MsgStore{
		Source: m,
	}

	go func() {
		for {
			fmt.Println("Store waiting for incoming data")
			select {
			case msg := <-dispatcher.StoreQ:
				fmt.Println("Store got some incoming data")
				store.Store(msg)
			}
		}
	}()

	return store
}

func (s *MsgStore) Store(msg *Msg) error {
	log.Printf("Store: %+v\n", msg)
	return nil
}
