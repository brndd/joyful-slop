package mappingrules

import "github.com/holoplot/go-evdev"

type MappingRule interface {
	MatchEvent(*evdev.InputDevice, *evdev.InputEvent) *evdev.InputEvent
}

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type SimpleMappingRule struct {
	Input  RuleTarget
	Output RuleTarget
}

// A Combo Mapping Rule can require multiple physical button presses for a single output button
type ComboMappingRule struct {
	Inputs []RuleTarget
	Output RuleTarget
	State  int
}

type RuleTarget struct {
	Device   *evdev.InputDevice
	Type     evdev.EvType
	Code     evdev.EvCode
	Inverted bool
}
