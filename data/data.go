package data

import "time"

// Data is an array of timestamps and values representing the same
// source of data over a period of time
type Data struct {
	Value     any           `json:"value"`
	TimeStamp time.Duration `json:"time"`
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
