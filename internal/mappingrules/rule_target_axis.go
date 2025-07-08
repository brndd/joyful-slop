package mappingrules

import (
	"github.com/holoplot/go-evdev"
)

type RuleTargetAxis struct {
	DeviceName    string
	Device        *evdev.InputDevice
	Axis          evdev.EvCode
	Inverted      bool
	DeadzoneStart int32
	DeadzoneEnd   int32
	Sensitivity   float64
}

func NewRuleTargetAxis(device_name string,
	device *evdev.InputDevice,
	axis evdev.EvCode,
	inverted bool,
	deadzone_start int32,
	deadzone_end int32,
	sensitivity float64) *RuleTargetAxis {

	return &RuleTargetAxis{
		DeviceName:    device_name,
		Device:        device,
		Axis:          axis,
		Inverted:      inverted,
		DeadzoneStart: deadzone_start,
		DeadzoneEnd:   deadzone_end,
		Sensitivity:   sensitivity,
	}
}

// TODO: lots of fixes and decisions to make here. Should we normalize all axes to the same range?
// How do we handle deadzones in light of that?
func (target *RuleTargetAxis) NormalizeValue(value int32) int32 {
	if !target.Inverted {
		return value
	}

	axisRange := target.DeadzoneEnd - target.DeadzoneStart
	axisMid := target.DeadzoneEnd - axisRange/2
	delta := value - axisMid
	if delta < 0 {
		delta = -delta
	}

	if value < axisMid {
		return axisMid + delta
	} else if value > axisMid {
		return axisMid - delta
	}

	// If we reach here, we're either exactly at the midpoint or something
	// strange has happened. Either way, just return the value.
	return value
}

func (target *RuleTargetAxis) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	// TODO: we can use the axis begin/end to decide whether to emit the event
	// TODO: oh no we need center deadzones actually...
	return &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  target.Axis,
		Value: value,
	}
}

func (target *RuleTargetAxis) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent) bool {
	return device == target.Device &&
		event.Type == evdev.EV_ABS &&
		event.Code == target.Axis
}
