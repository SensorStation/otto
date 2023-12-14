package main

import (
	"fmt"
	"time"
)

type Timestamp struct {
	Val		float64		`json:"val"`
	Time	int64		`json:"time"`
}

func NewTimestamp(v float64) (ts *Timestamp) {
	ts = &Timestamp{
		Val: v,
	}
	now := time.Now()
	ts.Time = now.Unix()
	return ts
}

func (ts *Timestamp) String() (str string) {
	t := time.Unix(ts.Time, 0)
	str = fmt.Sprintf("%s: %3.2f", t.Format(time.RFC3339), ts.Val)
	return str
}

type Timeseries struct {
	Values []*Timestamp		`json:"values"`
}

// NewTimeseries wraps a Timestap array with the given ID
func NewTimeseries() (ts *Timeseries) {
	ts = &Timeseries{}
	ts.Values = []*Timestamp{}
	return ts
}

// Add a value to the timeseries. A timestamp will be generated
// for the given piece of data
func (ts *Timeseries) Append(v float64) {
	t := NewTimestamp(v)
	ts.Values = append(ts.Values, t)
}

// Get a value from the Timeseries
func (ts *Timeseries) Get(idx int) *Timestamp {
	return ts.Values[idx]
}

func (ts *Timeseries) String() (str string) {
	for _, tts := range ts.Values {
		str += tts.String()
	}
	return str
}
