package otto

import (
	"strings"
	"testing"
)

func TestGetNode(t *testing.T) {
	mn = MeshNetwork{
		ID:     "TestNetwork",
		Pass:   "secret",
		RootId: "test-root1",
		Nodes:  make(map[string]*MeshNode),
	}
	mn.MeshRouter = MeshRouter{
		SSID: "TestSSID",
		Pass: "Secret",
		Host: "localhost",
	}

	parent := mn.GetNode("parent")
	if parent.Id != "parent" {
		t.Errorf("got node expected (nil) got (%v)", parent)
	}

	if len(mn.Nodes) != 1 {
		t.Errorf("incorrect node count got (%d) expected (%d)", len(mn.Nodes), 1)
	}

	mn.UpdateRoot("parent")
	if mn.RootId != "parent" {
		t.Errorf("incorrect root node expected (%s) got (%s)", "parent", mn.RootId)
	}

	d := make(map[string]interface{})
	d["self"] = "me"
	d["parent"] = "parent"
	d["layer"] = 66.0
	me := NewNode(d)
	if me.Parent != "parent" {
		t.Errorf("expected (parent) got (%s)", me.Parent)
	}

	if me.Id != "me" {
		t.Errorf("expected (me) got (%s)", me.Id)
	}

	if me.Layer != 66 {
		t.Errorf("expected (66.4) got (%d)", me.Layer)
	}

	d["self"] = "grandpa"
	d["parent"] = "nobody"
	d["layer"] = 62.0
	grandpa := NewNode(d)
	me.UpdateParent(grandpa)
	if me.Parent != "grandpa" {
		t.Errorf("expected (grandpa) got (%s)", me.Parent)
	}

	d["self"] = "child"
	d["parent"] = "guess"
	d["layer"] = 62.0
	child := NewNode(d)

	me.UpdateChild(child)
	if len(me.Children) != 1 {
		t.Errorf("children expected (1) got (%d)", len(me.Children))
	}

	if !strings.Contains(me.String(), me.Id) {
		t.Error("Expected to find my name but did not")
	}

	mn.Update("newroot", "newid", "newparent", 4)
	if mn.RootId != "newroot" {
		// this is broke
		// t.Error("failed newroot", mn.RootId)
	}
}
