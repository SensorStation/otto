package otto

type Store struct {
	Source map[string]map[string]float64
	StoreQ chan *Msg
}

func NewStore() *Store {
	m := make(map[string]map[string]float64)
	q := make(chan *Msg)
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

func (s *Store) Store(msg *Msg) error {
	l.Info("Store: ", "message", msg)
	return nil
}
