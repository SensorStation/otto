package otto

import "time"

// Timeseries represents a single source of data over a time period
type Timeseries struct {
	Label string
	Data  []*Data
}

// NewTimeseries will start a new data timeseries with the given label
func NewTimeseries(label string) *Timeseries {
	return &Timeseries{
		Label: label,
	}
}

// Add a new Data point to the given Timeseries
func (ts *Timeseries) Add(d any) *Data {
	dat := NewData(d)
	ts.Data = append(ts.Data, dat)
	return dat
}

// Len returns the number of data points contained in this timeseries
func (ts *Timeseries) Len() int {
	return len(ts.Data)
}

// Data is an array of timestamps and values representing the same
// source of data over a period of time
type Data struct {
	Value any
	time.Time
}

// NewData create a new peice data point with the current time
// as the timestamp
func NewData(v interface{}) *Data {
	return &Data{
		Value: v,
		Time:  time.Now(),
	}
}

// Return the float64 representation of the data. If the data is not
// represented by a float64 value a panic will follow
func (d *Data) Float() float64 {
	return d.Value.(float64)
}

// Int returns the integer value of the data. If the data is not
// an integer a panic will result.
func (d *Data) Int() int {
	return d.Value.(int)
}
