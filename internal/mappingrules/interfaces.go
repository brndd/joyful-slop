package mappingrules

import "github.com/holoplot/go-evdev"

type MappingRule interface {
	MatchEvent(*evdev.InputDevice, *evdev.InputEvent, *string) (*evdev.InputDevice, *evdev.InputEvent)
}

// RuleTargets represent either a device input to match on, or an output to produce.
// Some RuleTarget types may work via side effects, such as RuleTargetModeSelect.
type RuleTarget interface {
	// NormalizeValue takes the raw input value and possibly modifies it based on the Target settings.
	// (e.g., inverting the value if Inverted == true)
	NormalizeValue(int32) int32

	// CreateEvent typically takes the (probably normalized) value and returns an event that can be emitted
	// on a virtual device.
	//
	// For RuleTargetModeSelect, this method modifies the active mode and returns nil.
	//
	// TODO: should we normalize inside this function to simplify the interface?
	CreateEvent(int32, *string) *evdev.InputEvent
}
