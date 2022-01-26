package answer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var eventsFile = "events.json"

type EventType string

const (
	Create EventType = "create"
	Update EventType = "update"
	Delete EventType = "delete"
)

type Event struct {
	Event EventType `json:"event"`
	Data  Answer    `json:"data"`
}

type MapOfEvents map[string][]Event

func (allEvents MapOfEvents) CreateEvent(eventType EventType, data Answer) error {
	events, ok := allEvents[data.Key]
	if !ok {
		events = make([]Event, 0)
	}
	events = append(events, Event{Event: eventType, Data: data})

	allEvents[data.Key] = events
	err := allEvents.StoreAllEvents() // update all events
	if err != nil {
		return err
	}
	return nil
}

func (allEvents MapOfEvents) StoreAllEvents() error {
	file, err := json.MarshalIndent(allEvents, "", "")
	if err != nil {
		return fmt.Errorf("storeAllEvents: %v", err)
	}
	err = ioutil.WriteFile(eventsFile, file, 0644)
	if err != nil {
		return fmt.Errorf("storeAllEvents: %v", err)
	}
	return nil
}

// InitEvents will populate internal events map using the json events file
// If file is not found or it is not valid, an empty map will be created instead
func InitEvents() MapOfEvents {
	allEventsIndexed := make(MapOfEvents)

	body, err := ioutil.ReadFile(eventsFile)
	if err != nil {
		return allEventsIndexed
	}

	err = json.Unmarshal(body, &allEventsIndexed)
	if err != nil {
		return allEventsIndexed
	}

	return allEventsIndexed
}
