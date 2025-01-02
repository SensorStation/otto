package devices

import (
	"github.com/tarm/serial"
)

type SerialDevice struct {
	PortName string
	Baud     int
	*serial.Port
	*Device
}

func NewSerialDevice(name, port string, baud int) *SerialDevice {
	sd := &SerialDevice{PortName: port, Baud: baud}
	sd.Device = NewDevice(name)
	return sd
}

func (s *SerialDevice) Open() (err error) {
	c := &serial.Config{Name: s.PortName, Baud: s.Baud}
	s.Port, err = serial.OpenPort(c)
	if err != nil {
		return err
	}
	return nil
}

func (s *SerialDevice) Write(buf []byte) (int, error) {
	return s.Port.Write(buf)
}

func (s *SerialDevice) Read(buf []byte) (int, error) {
	return s.Port.Read(buf)
}
