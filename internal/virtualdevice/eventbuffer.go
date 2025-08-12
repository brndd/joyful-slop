// Code to manage sending events on the Virtual Device
package virtualdevice

import (
	"syscall"
	"time"

	"github.com/holoplot/go-evdev"
)

type EventBuffer struct {
	events []*evdev.InputEvent
	Device VirtualDevice
	Name   string
}

func (buffer *EventBuffer) AddEvent(event *evdev.InputEvent) {
	buffer.events = append(buffer.events, event)
}

func (buffer *EventBuffer) SendEvents() []error {
	eventTime := syscall.NsecToTimeval(int64(time.Now().Nanosecond()))
	writeErrors := make([]error, 0)

	if len(buffer.events) == 0 {
		return writeErrors
	}

	for i := 0; i < len(buffer.events); i++ {
		buffer.events[i].Time = eventTime
		err := buffer.Device.WriteOne(buffer.events[i])
		if err != nil {
			writeErrors = append(writeErrors, err)
		}
	}

	err := buffer.Device.WriteOne(&evdev.InputEvent{
		Time:  eventTime,
		Type:  evdev.EV_SYN,
		Code:  evdev.SYN_REPORT,
		Value: 0,
	})

	if err != nil {
		writeErrors = append(writeErrors, err)
	}

	buffer.events = make([]*evdev.InputEvent, 0, 100)
	return writeErrors
}
