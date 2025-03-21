package data

import (
	"testing"

	"github.com/sensorstation/otto/messanger"
)

func TestTopics(t *testing.T) {
	tctl := "ss/c/station/test"
	topic := messanger.TopicControl("test")
	if topic != tctl {
		t.Errorf("expected topic (%s) got (%s)", tctl, topic)
	}
	tdat := "ss/d/station/test"
	topic = messanger.TopicData("test")
	if topic != tdat {
		t.Errorf("expected topic (%s) got (%s)", tdat, topic)
	}

	messanger.TopicData("test")
	messanger.TopicData("test")

	var v int
	var ex bool
	if v, ex = messanger.Topics[tctl]; !ex {
		t.Errorf("Expected to find %s but did not", tctl)
	}
	if v != 1 {
		t.Errorf("Expected (1) got (%d)", v)
	}

	if v, ex = messanger.Topics[tdat]; !ex {
		t.Errorf("Expected (3) got (%d)", v)
	}

}
