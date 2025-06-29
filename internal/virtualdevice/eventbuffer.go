// Code to manage sending events on the Virtual Device
package virtualdevice

import (
	"syscall"
	"time"

	"github.com/holoplot/go-evdev"
)

type EventBuffer struct {
	events []*evdev.InputEvent
	Device *evdev.InputDevice
}

func NewEventBuffer(device *evdev.InputDevice) *EventBuffer {
	return &EventBuffer{
		events: make([]*evdev.InputEvent, 0),
		Device: device,
	}
}

func (buffer *EventBuffer) AddEvent(event *evdev.InputEvent) {
	buffer.events = append(buffer.events, event)
}

func (buffer *EventBuffer) SendEvents() {
	eventTime := syscall.NsecToTimeval(int64(time.Now().Nanosecond()))

	for i := 0; i < len(buffer.events); i++ {
		buffer.events[i].Time = eventTime
		buffer.Device.WriteOne(buffer.events[i])
	}

	buffer.Device.WriteOne(&evdev.InputEvent{
		Time:  eventTime,
		Type:  evdev.EV_SYN,
		Code:  evdev.SYN_REPORT,
		Value: 0,
	})

	buffer.events = make([]*evdev.InputEvent, 0)
}
