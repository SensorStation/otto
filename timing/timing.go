package timing

import "time"

var StartTime time.Time

func init() {
	StartTime = time.Now()
}

func Timestamp() time.Duration {
	return time.Since(StartTime)
}
