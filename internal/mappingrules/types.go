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

type RuleTarget interface {
	NormalizeValue(int32) int32
	CreateEvent(int32, *string) *evdev.InputEvent
	GetCode() evdev.EvCode
	GetDeviceName() string
	GetDevice() *evdev.InputDevice
}

type RuleTargetBase struct {
	DeviceName string
	Device     *evdev.InputDevice
	Code       evdev.EvCode
	Inverted   bool
}

type RuleTargetButton struct {
	RuleTargetBase
}

type RuleTargetAxis struct {
	RuleTargetBase
	AxisStart   int32
	AxisEnd     int32
	Sensitivity float64
}

type RuleTargetModeSelect struct {
	RuleTargetBase
	ModeSelect []string
}
