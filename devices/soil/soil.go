package main

import (
	"github.com/sensorstation/otto/device"
	"golang.org/x/exp/io/i2c"
)

type Soil struct {
	*device.Device
}

func New(name, bus string, addr int) *Soil {
	s := &Soil{}
	d, err := i2c.Open(&i2c.Devfs{Dev: bus}, addr)
	if err != nil {
		panic(err)
	}
	defer d.Close()

	// for {
	// 	wbuf := []byte{0x0f, 0x10}
	// 	err := d.Write(wbuf)
	// 	if err != nil {
	// 		fmt.Printf("error: %+v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	time.Sleep(3 * time.Millisecond)

	// 	var rbuf []byte = []byte{0x0, 0x0}
	// 	err = d.Read(rbuf)
	// 	if err != nil {
	// 		fmt.Printf("error: %+v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	val := rbuf[0]<<8 | rbuf[1]
	// 	fmt.Printf("moisture: %v => %d\n", rbuf, val)
	// 	time.Sleep(time.Second)
	// 	continue

	// 	// get temp
	// 	wbuf = []byte{0x00, 0x04}
	// 	err = d.Write(wbuf)
	// 	if err != nil {
	// 		fmt.Printf("error: %+v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	time.Sleep(3 * time.Millisecond)
	// 	rbuf = make([]byte, 4)
	// 	err = d.Read(rbuf)
	// 	if err != nil {
	// 		fmt.Printf("error: %+v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	val = rbuf[0]<<24 | rbuf[1]<<16 | rbuf[2]<<8 | rbuf[3]
	// 	fmt.Printf("nval: %T, %v\n", val, val)
	// 	fval := float64(val)
	// 	fmt.Printf("fval: %T, %5.2f\n", fval, fval)

	// 	tc := 0.000015259 * fval
	// 	fmt.Printf("  tc: %T, %v\n", tc, tc)

	// 	tf := (tc * (9.0 / 5.0)) + 32.0
	// 	fmt.Printf("tempc: %v => %5.2f, tempf: %5.2f\n", rbuf, tc, tf)

	// 	time.Sleep(time.Second)
	// }
	return s
}

// func readSoilMoisture(bus i2c.BusCloser, addr uint16) (uint16, error) {
// 	// Assuming the sensor uses a specific register for moisture reading
// 	// Replace 0x00 with the correct register address
// 	reg := []byte{0x00}
// 	data := make([]byte, 2)

// 	if err := bus.Tx(addr, reg, data); err != nil {
// 		return 0, err
// 	}

// 	// Combine the two bytes to get the moisture value
// 	moisture := (uint16(data[0]) << 8) | uint16(data[1])

// 	return moisture, nil
// }

func main() {
	panic("not implemented")
}
