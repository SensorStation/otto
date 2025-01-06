package data

import (
	"log/slog"
	"os"

	"github.com/sensorstation/otto/message"
)

type Store struct {
	Filename string
	Source   map[string]map[string]float64
	StoreQ   chan *message.Msg

	f *os.File
}

var (
	store *Store
)

func GetStore() *Store {
	if store == nil {
		store = NewStore()
	}
	return store
}

func NewStore() *Store {
	m := make(map[string]map[string]float64)
	q := make(chan *message.Msg)
	store := &Store{
		Source: m,
		StoreQ: q,
	}

	go func() {
		for {
			select {
			case msg := <-store.StoreQ:
				store.Store(msg)
			}
		}
	}()

	return store
}

func (s *Store) Store(msg *message.Msg) error {
	slog.Info("Store: ", "message", msg)
	return nil
}
