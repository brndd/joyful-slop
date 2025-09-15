package mappingrules

import (
	"time"

	"github.com/holoplot/go-evdev"
)

type MappingRule interface {
	MatchEvent(Device, *evdev.InputEvent, *string) (*evdev.InputDevice, *evdev.InputEvent)
}

type TimedEventEmitter interface {
	TimerEvent() *evdev.InputEvent
	GetOutputDevice() *evdev.InputDevice
}

// RuleTargets represent either a device input to match on, or an output to produce.
// Some RuleTarget types may work via side effects, such as RuleTargetModeSelect.
type RuleTarget interface {
	// NormalizeValue takes the raw input value and possibly modifies it based on the Target settings.
	// (e.g., inverting the value if Inverted == true)
	NormalizeValue(int32) int32

	// CreateEvent creates an event that can be emitted on a virtual device.
	// For RuleTargetModeSelect, this method modifies the active mode and returns nil.
	//
	// CreateEvent does not do any error checking; it assumes it is receiving verified and sanitized output.
	// The caller is responsible for determining whether an event *should* be emitted.
	//
	// Typically int32 is the input event's normalized value. *string is the current mode, but is optional
	// for most implementations.
	CreateEvent(int32, *string) *evdev.InputEvent

	// MatchEvent returns true if the provided device and input event are a match for this rule target
	MatchEvent(device Device, event *evdev.InputEvent) bool
}

// Device is an interface abstraction on top of evdev.InputDevice, implementing
// only the methods we need in this package. This is used for testing, and the
// Device can be safely cast to an *evdev.InputDevice when necessary.
type Device interface {
	AbsInfos() (map[evdev.EvCode]evdev.AbsInfo, error)
}

const (
	AxisValueMin = int32(-32768)
	AxisValueMax = int32(32767)
	NoNextEvent  = time.Duration(-1)
)
