package main

import "github.com/holoplot/go-evdev"

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
}
