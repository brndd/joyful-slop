package mappingrules

import (
	"github.com/holoplot/go-evdev"
)

type RuleTargetRelaxis struct {
	DeviceName string
	Device     Device
	Axis       evdev.EvCode
	Inverted   bool
}

func NewRuleTargetRelaxis(device_name string,
	device Device,
	axis evdev.EvCode,
	inverted bool) (*RuleTargetRelaxis, error) {

	return &RuleTargetRelaxis{
		DeviceName: device_name,
		Device:     device,
		Axis:       axis,
		Inverted:   inverted,
	}, nil
}

// NormalizeValue takes a raw input value and converts it to a value suitable for output.
//
// Relative axes are currently only supported for output.
// TODO: make this have an error return?
func (target *RuleTargetRelaxis) NormalizeValue(value int32) int32 {
	return 0
}

func (target *RuleTargetRelaxis) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	return &evdev.InputEvent{
		Type:  evdev.EV_REL,
		Code:  target.Axis,
		Value: value,
	}
}

// Relative axis is only supported for output.
func (target *RuleTargetRelaxis) MatchEvent(device Device, event *evdev.InputEvent) bool {
	return false
}
