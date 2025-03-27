package messanger

import (
	"encoding/json"
	"net/http"
	"time"
)

var (
	events *EventLog
)

// Event is anything significant that we would like to track such as a
// switch going off or a threshold being surpassed
type Event struct {
	Name string
	Data any
	time.Time
}

// EventLog will keep track of a series of events that we are
// interested in tracking
type EventLog struct {
	Events    []*Event
	MaxEvents int
}

// Get the global EventLog
func GetEventLog() *EventLog {
	if events == nil {
		events = &EventLog{}
	}
	return events
}

// AddEvent takes the data as provided, creates an event from it then
// inserts the event into the event log
func AddEvent(name string, t time.Time, d any) {
	event := &Event{
		Name: name,
		Time: t,
		Data: d,
	}
	events := GetEventLog()
	events.Add(event)
}

// AddEventNow adds the given event at the current time, that is the
// time the call was made
func AddEventNow(name string, d any) {
	AddEvent(name, time.Now(), d)
}

// Add an Event to the EventLog
func (e *EventLog) Add(event *Event) {
	e.Events = append(e.Events, event)
}

// ServeHTTP will return the exisint EventLog, in the future we'll add
// the ability to specify how many events we want to retrieve.
func (evl *EventLog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(evl).Encode(ms)
}
