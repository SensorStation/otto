package main

import (
	"testing"
	"time"

	"github.com/sensorstation/otto/device/drivers"
	"github.com/sensorstation/otto/messanger"
)

func TestBlink(t *testing.T) {
	drivers.GetGPIO().Mock = true
	messanger.SetMQTTClient(messanger.GetMockClient())

	led, done := initLED("test-led", 13)
	if led.Device.Name() != "test-led" {
		t.Errorf("name expected (%s) got (%s)", "test-led", led.Device.Name())
	}

	domqtt(led)
	go dotimer(led, 100*time.Millisecond, done)
	time.Sleep(1 * time.Second)
	done <- true

	if count != 10 {
		t.Errorf("count expected (%d) got (%d)", 10, count)
	}
}
