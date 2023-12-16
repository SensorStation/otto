package iote

import (
	"fmt"
)

var (
	Store *Storage
)

type Storage struct {
	Source map[string]map[string]float64
}

func init() {
	Store = NewStore()
}

func NewStore() *Storage {
	m := make(map[string]map[string]float64)
	return &Storage{
		Source: m,
	}
}

func (s *Storage) Store(msg *Msg) {
	fmt.Printf("MSG: %+v\n", msg)
}
