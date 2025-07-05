package main

import (
	"time"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

const (
	TimerCheckIntervalMs  = 250
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

func timerWatcher(rule *mappingrules.MappingRuleProportionalAxis, channel chan<- ChannelEvent) {
	for {
		event := rule.TimerEvent()
		if event != nil {
			channel <- ChannelEvent{
				Device: rule.Output.(*mappingrules.RuleTargetModeSelect).Device,
				Event:  event,
				Type:   ChannelEventTimer,
			}
		}
		time.Sleep(TimerCheckIntervalMs * time.Millisecond)
	}
}
