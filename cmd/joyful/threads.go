package main

import (
	"time"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

const (
	TimerCheckIntervalMs  = 1
	DeviceCheckIntervalMs = 1
)

func eventWatcher(device *evdev.InputDevice, channel chan<- ChannelEvent) {
	for {
		event, err := device.ReadOne()
		if err != nil {
			logger.LogError(err, "Error while reading event. Disconnecting device.")
			return
		}
		channel <- ChannelEvent{Device: device, Event: event, Type: ChannelEventInput}

		if event.Type == evdev.EV_SYN {
			time.Sleep(DeviceCheckIntervalMs * time.Millisecond)
		}
	}
}

func timerWatcher(rule mappingrules.TimedEventEmitter, channel chan<- ChannelEvent) {
	for {
		event := rule.TimerEvent()
		if event != nil {
			channel <- ChannelEvent{
				Device: rule.GetOutputDevice(),
				Event:  event,
				Type:   ChannelEventTimer,
			}
		}
		time.Sleep(TimerCheckIntervalMs * time.Millisecond)
	}
}
