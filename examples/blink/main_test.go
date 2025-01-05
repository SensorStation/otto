package main

import (
	"testing"
	"time"

	"github.com/sensorstation/otto/devices"
	"github.com/sensorstation/otto/messanger"
)

func TestBlink(t *testing.T) {
	devices.GetGPIO().Mock = true
	messanger.SetMQTTClient(messanger.GetMockClient())

	led, done := initLED("test-led", 13)
	if led.Name != "test-led" {
		t.Errorf("name expected (%s) got (%s)", "test-led", led.Name)
	}

	if led.Pub != messanger.TopicData(led.Name) {
		t.Errorf("publish topic expected (%s) got (%s)", led.Pub, messanger.TopicData(led.Name))
	}

	if led.Subs != nil {
		t.Errorf("expected sub to be nil but got (%q)", led.Subs)
	}

	if led.EvtQ != nil {
		t.Errorf("expected led EvtQ to be nil but it is not")
	}

	if led.Period != time.Duration(0) {
		t.Errorf("expected led period to be zero but got (%d)", led.Period)
	}

	domqtt(led)
	if len(led.Subs) != 1 {
		t.Errorf("subs expected (%d) got (%q)", 1, len(led.Subs))
	}

	led.Period = 100 * time.Millisecond
	go dotimer(led, done)
	time.Sleep(1 * time.Second)
	done <- true

	if led.Period != 100*time.Millisecond {
		t.Errorf("count expected (%d ms) got (%d)", 10, count)
	}

	if count != 10 {
		t.Errorf("count expected (%d) got (%d)", 10, count)
	}

	_ = done
}
