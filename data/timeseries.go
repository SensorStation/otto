package data

import "time"

// Timeseries represents a single source of data over a time period
type Timeseries struct {
	Station string    `json:"station"`
	Label   string    `json:"label"`
	Data    []*Data   `json:"data"`
	Start   time.Time `json:"start"`
}

// NewTimeseries will start a new data timeseries with the given label
func NewTimeseries(station, label string) *Timeseries {
	return &Timeseries{
		Station: station,
		Label:   label,
		Start:   time.Now(),
	}
}

// Add a new Data point to the given Timeseries
func (ts *Timeseries) Add(d any) *Data {
	dat := &Data{
		Value:     d,
		TimeStamp: time.Since(ts.Start),
	}
	ts.Data = append(ts.Data, dat)
	return dat
}

// Len returns the number of data points contained in this timeseries
func (ts *Timeseries) Len() int {
	return len(ts.Data)
}
