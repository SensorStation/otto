package messanger

import "testing"

func TestMessanger(t *testing.T) {
	id := "testID"
	topic := "testTopic"
	m := NewMessanger(id, topic)
	if m.ID != id {
		t.Errorf("expected ID (%s) got (%s)", id, m.ID)
	}

	if m.Topic != topic {
		t.Errorf("expected topic (%s) got (%s)", topic, m.Topic)
	}

	if m.Published != 0 {
		t.Errorf("expected published (%d) got (%d)", 0, m.Published)
	}

	if m.PubSub != mqtt {
		t.Errorf("expected PubSub (mqtt) got (%p)", mqtt)
	}

	if len(m.Subs) > 1 {
		t.Errorf("expected len m.Subs (0) got (%d)", len(m.Subs))
	}
}
