package iote

import (
	"fmt"
	"math/rand"
	"net/http"
)

// PeriodicRandomData will collected a new random piece of data
// every period and transmit it to the given mqtt channel
type Rando struct {
	F float64
}

func NewRando() (r *Rando) {
	r = &Rando{
		F: 0.0,
	}
	return r
}

func (p Rando) Get() interface{} {
	p.F = rand.Float64()
	s := fmt.Sprintf("%f", p.F)
	return s
}

func (p Rando) GetFloat() float64 {
	return rand.Float64()
}

func (p Rando) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Get()
	fmt.Fprintf(w, "%f", p.F)
}
