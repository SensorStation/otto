package device

import "github.com/tarm/serial"

type SerialPort struct {
	portName string
	Baud     int
	*serial.Port
}

func GetSerialPort(name, port string, opts any) *SerialPort {

	var baud int
	switch opts.(type) {
	case int:
		baud = opts.(int)
	}

	sd := &SerialPort{portName: port, Baud: baud}
	return sd
}

func (s *SerialPort) Open() (err error) {
	c := &serial.Config{Name: s.portName, Baud: s.Baud}
	s.Port, err = serial.OpenPort(c)
	if err != nil {
		return err
	}
	return nil
}

func (d *SerialPort) PortName() string {
	return d.portName
}

func (s SerialPort) Write(buf []byte) (int, error) {
	return s.Port.Write(buf)
}

func (s SerialPort) Read(buf []byte) (int, error) {
	return s.Port.Read(buf)
}
