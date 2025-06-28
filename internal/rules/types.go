package rules

import "github.com/holoplot/go-evdev"

type RuleTarget struct {
	Device   *evdev.InputDevice
	Type     evdev.EvType
	Code     evdev.EvCode
	Inverted bool
}
