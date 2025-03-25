package messanger

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Topics maintains the list of topics used by otto and the
// applications. It maintains the topic format and a count for each
// time the topic is used
type Topics struct {
	StationName string
	TopicFmt    string
	Topicmap    map[string]int
}

var (
	topics *Topics
)

func init() {
	topics = &Topics{
		TopicFmt:    "ss/%s/%s/%s",
		Topicmap:    make(map[string]int),
		StationName: "",
	}
}

// GetTopics will return the Topics structure, one per application.
func GetTopics() *Topics {
	return topics
}

// SetStationName will use the value to set the station that will be
// used when publishing messages from this station.
func (t *Topics) SetStationName(name string) {
	t.StationName = name
}

// Control will return a control topic e.g. ss/c/station/foo
func (t *Topics) Control(topic string) string {
	top := fmt.Sprintf(t.TopicFmt, "c", t.StationName, topic)
	t.Topicmap[top]++
	return top
}

// Control will return a data topic e.g. ss/d/station/foo
func (t *Topics) Data(topic string) string {
	top := fmt.Sprintf(t.TopicFmt, "d", t.StationName, topic)
	t.Topicmap[top]++
	return top
}

// ServeHTTP is a JSON endpoint that returns all the topics used by
// this station.
func (t Topics) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jstr, err := json.Marshal(t)
	if err != nil {
		http.Error(w, "Not Yet Supported", 401)
	}
	fmt.Fprint(w, jstr)
}
