package device

// import (
// 	"log"
// 	"time"
// )

// func NewAnalogDevice(name string, offset int, opts ...any) *Device {
// 	d := NewDevice(name)
// 	a := GetADS1115()
// 	p, err := a.Pin(name, offset, opts)
// 	if err == nil {
// 		d.ReadWriteCloser = p
// 	}
// 	d.Error = err
// 	return d

// }

// type AnalogMock struct {
// }

// func GetAnalogMock(name string) *AnalogMock {
// 	return &AnalogMock{
// 		BaseDevice: NewDevice(name),
// 	}
// }

// func (d *AnalogMock) Get() (val float64, err error) {
// 	return val, err
// }

// func (d *AnalogMock) Set(val float64) error {
// 	return nil
// }

// func (d *AnalogMock) Close() error {
// 	return nil
// }

// func (d *AnalogMock) ReadContinuous() <-chan float64 {
// 	q := make(chan float64)

// 	go func() {
// 		for {
// 			v, err := d.Get()
// 			if err != nil {
// 				log.Printf("%s ReadContinuous failed read %s", d.name, err)
// 				continue
// 			}
// 			q <- v
// 			<-time.After(5 * time.Second)
// 		}
// 	}()

// 	return q
// }
