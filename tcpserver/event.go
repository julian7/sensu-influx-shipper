package tcpserver

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/sensu/sensu-go/types"
)

type Event struct {
	*types.Event
}

func ReadEvent(reader io.Reader) (*Event, error) {
	eventRaw, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, NewError("read", err)
	}

	sensuEvent := &types.Event{}
	if err := json.Unmarshal(eventRaw, sensuEvent); err != nil {
		return nil, NewError("json parsing", err)
	}

	event := &Event{Event: sensuEvent}

	if err := event.Validate(); err != nil {
		return nil, NewError("validate", err)
	}

	return event, nil
}

func (e *Event) Validate() error {
	if e.Event.Timestamp <= 0 {
		return errors.New("timestamp is missing or must be greater than zero")
	}

	return e.Event.Validate()
}

func (e *Event) BaseTags() map[string]string {
	var ipaddr string

	for _, intf := range e.Event.Entity.System.Network.Interfaces {
		if len(intf.MAC) > 0 && len(intf.Addresses) > 0 {
			ipaddr = intf.Addresses[0]
			break
		}
	}

	tags := map[string]string{}

	for _, item := range []struct {
		key   string
		value string
	}{
		{"host", e.Event.Entity.Name},
		{"ip", ipaddr},
		{"check", e.Event.Check.Name},
	} {
		if _, ok := tags[item.key]; !ok {
			tags[item.key] = item.value
		}
	}

	return tags
}
