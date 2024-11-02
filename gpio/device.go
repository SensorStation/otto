package main

// import "fmt"

// var (
// 	devices Devices
// )

// type Devices map[string]*Device

// func NewDevice(name string) *Device {
// 	d := &Device{}
// 	d.Pins = make(map[int]*Pin)

// 	devices[name] = d
// 	return d
// }

// func (dev *Device) AddPin(pin *Pin) {
// 	dev.Pins[pin.Offset] = pin
// }

// func (dev *Device) String() string {
// 	output := dev.Name

// 	for n, p := range dev.Pins {
// 		output += fmt.Sprintf("%d, %s\n", n, p.String())
// 	}

// 	return output
// }

// func (d Devices) String() string {
// 	output := ""
// 	for n, d := range devices {
// 		output += fmt.Sprintf("%s: %s\n", n, d.String())
// 	}

// 	return output
// }
