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

func NewData(v interface{}) *Data {
	return &Data{
		Value: v,
		Time:  time.Now(),
	}
}

func (d *Data) Float() float64 {
	return d.Value.(float64)
}

func NewTimeseries(label string) *Timeseries {
	return &Timeseries{
		Label: label,
	}
}
