package mappingrules

import (
	"github.com/holoplot/go-evdev"
)

type RuleTargetAxis struct {
	RuleTargetBase
	AxisStart   int32
	AxisEnd     int32
	Sensitivity float64
}

func NewRuleTargetAxis(device_name string,
	device *evdev.InputDevice,
	code evdev.EvCode,
	inverted bool,
	axis_start int32,
	axis_end int32,
	sensitivity float64) *RuleTargetAxis {

	return &RuleTargetAxis{
		RuleTargetBase: NewRuleTargetBase(device_name, device, code, inverted),
		AxisStart:      axis_start,
		AxisEnd:        axis_end,
		Sensitivity:    sensitivity,
	}
}

func (target *RuleTargetAxis) NormalizeValue(value int32) int32 {
	if !target.Inverted {
		return value
	}

	axisRange := target.AxisEnd - target.AxisStart
	axisMid := target.AxisEnd - axisRange/2
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
		Code:  target.Code,
		Value: value,
	}
}
