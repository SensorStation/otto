package otto

import "github.com/sensorstation/otto/message"

type Store struct {
	Source map[string]map[string]float64
	StoreQ chan *message.Msg
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
	l.Info("Store: ", "message", msg)
	return nil
}
