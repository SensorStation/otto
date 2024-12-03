package otto

import "time"

type Timeseries struct {
	Label string
	Data  []*Data
}

type Data struct {
	Value interface{}
	time.Time
}

func NewTimeseries(label string) *Timeseries {
	return &Timeseries{
		Label: label,
	}
}
