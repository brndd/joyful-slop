package mappingrules

import (
	"time"

	"github.com/holoplot/go-evdev"
)

type MappingRule interface {
	MatchEvent(*evdev.InputDevice, *evdev.InputEvent, *string) *evdev.InputEvent
	OutputName() string
}

type MappingRuleBase struct {
	Name   string
	Output RuleTarget
	Modes  []string
}

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type SimpleMappingRule struct {
	MappingRuleBase
	Input RuleTarget
}

// A Combo Mapping Rule can require multiple physical button presses for a single output button
type ComboMappingRule struct {
	MappingRuleBase
	Inputs []RuleTarget
	State  int
}

type LatchedMappingRule struct {
	MappingRuleBase
	Input RuleTarget
	State bool
}

// TODO: How are we going to implement this? It needs to operate on a timer...
type ProportionalAxisMappingRule struct {
	MappingRuleBase
	Input     RuleTarget
	Output    RuleTarget
	LastEvent time.Time
}

type RuleTarget struct {
	DeviceName string
	ModeSelect []string
	Device     *evdev.InputDevice
	Type       evdev.EvType
	Code       evdev.EvCode
	Inverted   bool
}
