package messanger

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTopics(t *testing.T) {
	topics := GetTopics()
	if topics == nil {
		t.Error("GetTopics() returns nil topic expected it to be ready to go ")
	}

	name := "test-station"
	topics.SetStationName(name)
	if topics.StationName != name {
		t.Errorf("Expected station name (%s) got (%s)", name, topics.StationName)
	}

	ctltopic := "ss/c/test-station/controller"
	ctl := topics.Control("controller")
	if ctl != ctltopic {
		t.Errorf("Expected control topic (%s) got (%s)", ctltopic, ctl)
	}

	datatopic := "ss/d/test-station/controller"
	dat := topics.Control("data")
	if ctl != ctltopic {
		t.Errorf("Expected data topic (%s) got (%s)", datatopic, dat)
	}
}

func TestTopicHTTP(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		topics := GetTopics()
		topics.ServeHTTP(w, r)
	}

	req := httptest.NewRequest("GET", "/api/topics", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Expected err (nil) got (%s)", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected response code (200) got (%d)", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected response Header (application/json) got (%s)", contentType)
	}

	var topics Topics
	err = json.Unmarshal(body, &topics)
	if err != nil {
		t.Errorf("Expected err nil got (%s)", err)
	}

	bodystr := fmt.Sprintf("%s", body)
	println(bodystr)

}
