package device

import "github.com/tarm/serial"

type Serial struct {
	portName string
	Baud     int
	*serial.Port
}

func GetSerial(name, port string, opts any) *Serial {

	var baud int
	switch opts.(type) {
	case int:
		baud = opts.(int)
	}

	sd := &Serial{portName: port, Baud: baud}
	return sd
}

func (s *Serial) Open() (err error) {
	c := &serial.Config{Name: s.portName, Baud: s.Baud}
	s.Port, err = serial.OpenPort(c)
	if err != nil {
		return err
	}
	return nil
}

func (d *Serial) PortName() string {
	return d.portName
}

func (s Serial) Write(buf []byte) (int, error) {
	return s.Port.Write(buf)
}

func (s Serial) Read(buf []byte) (int, error) {
	return s.Port.Read(buf)
}
