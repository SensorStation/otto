package iote

import (
	"fmt"
)

var (
	store *Store
)

func init() {
	store = NewStore()
}

type Store struct {
	Source map[string]map[string]float64
}

func NewStore() *Store {
	m := make(map[string]map[string]float64)
	return &Store{
		Source: m,
	}
}

func (s *Store) Store(msg *Msg) {
	fmt.Printf("MSG: %+v\n", msg)
}
