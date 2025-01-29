package bme280

import (
	"testing"
	"time"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
)

func init() {
	// Setup the mock mqtt client
	c := messanger.GetMockClient()
	messanger.SetMQTTClient(c)
}

func TestBME280(t *testing.T) {
	name := "bme-test"
	bus := "/dev/i2c-fake"
	addr := 0x76

	devices.Mock = true
	bme := New(name, bus, addr)
	if bme == nil {
		t.Error("Failed to create bme device")
	}

	if bme.Name() != name {
		t.Errorf("expected name (%s) got (%s)", name, bme.Name())
	}

	if bme.Bus != bus {
		t.Errorf("expected bus (%s) got (%s)", bus, bme.Bus)
	}
	if bme.Addr != addr {
		t.Errorf("expected addr (%x) got (%x)", addr, bme.Addr)
	}

	err := bme.Init()
	if err != nil {
		t.Error("failed to initialize mock BME device")
	}

	resp, err := bme.Read()
	if err != nil {
		t.Error("failed to read from mock BME device")
	}

	// XXX these values are hard coded in the bme280 device mock read
	// may change to random values or other
	if resp.Temperature != 20.33 {
		t.Errorf("Failed to read temperature")
	}
	if resp.Pressure != 1027.33 {
		t.Errorf("Failed to read temperature")
	}
	if resp.Humidity != 74.33 {
		t.Errorf("Failed to read temperature")
	}

	// Set up for bme EventLoop run the loop every 200 milliseconds then
	// stop the loop
	count := 0
	topic := messanger.TopicData(name)
	bme.AddPub(topic)
	bme.Subscribe(topic, func(msg *messanger.Msg) {
		if msg.Topic != topic {
			t.Errorf("expected topic (%s) got (%s)", topic, msg.Topic)
			return
		}
		mmm, err := msg.Map()
		if err != nil {
			t.Errorf("failed to map bme280 response %s", err)
			return
		}

		for key, val := range mmm {
			switch key {
			case "Temperature":
				if val != 20.33 {
					t.Errorf("%s expected (20.33) got (%f)", key, val)
					return
				}

			case "Humidity":
				if val != 74.33 {
					t.Errorf("%s expected (74.33) got (%4.2f)", key, val)
					return
				}

			case "Pressure":
				if val != 1027.33 {
					t.Errorf("%s expected (1027.33) got (%4.2f)", key, val)
					return
				}

			default:
				t.Errorf("Unexpected key value %s - %v", key, val)
				return
			}
		}
		count++
	})

	done := make(chan any)
	go bme.TimerLoop(100*time.Millisecond, done, bme.ReadPub)

	select {
	case <-done:
		break

	case <-time.After(1 * time.Second):
		done <- true
	}
	if count < 5 {
		t.Errorf("Expected to recieve messanges expected (>= 5) got (%d) ", count)
	}
}
