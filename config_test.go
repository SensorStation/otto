package iote

import (
	"testing"
)

func TestConfig(t *testing.T) {
	config := GetConfig()

	if config.Broker != "tcp://localhost:1883" {
		t.Errorf("Expected broker (tcp://localhost:1883) got (%s)", config.Broker)
	}
}
