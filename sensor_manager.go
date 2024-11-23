package otto

import "fmt"

type SensorManager struct {
	Sensors map[string]*Sensor
}

func NewSensorManager() (sm *SensorManager) {
	sm = &SensorManager{
		Sensors: make(map[string]*Sensor),
	}
	return sm
}

func (sm *SensorManager) Callback(msg *Msg) {
	fmt.Printf("sensor manager: %+v\n", msg)
}
