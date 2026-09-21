package main

import (
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

type ChannelEventType int

const (
	ChannelEventInput ChannelEventType = iota
	ChannelEventTimer
	ChannelEventReload
)

type ChannelEvent struct {
	Type   ChannelEventType
	Device *evdev.InputDevice
	Event  *evdev.InputEvent
	Rule   mappingrules.MainLoopTimedEventEmitter
}
