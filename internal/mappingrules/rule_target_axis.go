package mappingrules

import (
	"errors"
	"fmt"

	"github.com/holoplot/go-evdev"
)

type RuleTargetAxis struct {
	DeviceName    string
	Device        RuleTargetDevice
	Axis          evdev.EvCode
	Inverted      bool
	DeadzoneStart int32
	DeadzoneEnd   int32
	axisSize      int32
	deadzoneSize  int32
}

func NewRuleTargetAxis(device_name string,
	device RuleTargetDevice,
	axis evdev.EvCode,
	inverted bool,
	deadzoneStart int32,
	deadzoneEnd int32) (*RuleTargetAxis, error) {

	info, err := device.AbsInfos()

	if err != nil {
		// If we can't get AbsInfo (for example, we're a virtual device)
		// we set the bounds to the maximum allowable
		info = map[evdev.EvCode]evdev.AbsInfo{
			axis: {
				Minimum: AxisValueMin,
				Maximum: AxisValueMax,
			},
		}
	}

	if _, ok := info[axis]; !ok {
		return nil, fmt.Errorf("device does not support axis %v", axis)
	}

	if deadzoneStart > deadzoneEnd {
		return nil, errors.New("deadzone_end must be a higher value than deadzone_start")
	}

	deadzoneSize := Abs(deadzoneEnd - deadzoneStart)

	// Our output range is limited to 16 bits, but we represent values internally with 32 bits.
	// As a result, we shouldn't need to worry about integer overruns
	axisSize := info[axis].Maximum - info[axis].Minimum - deadzoneSize

	if axisSize == 0 {
		return nil, errors.New("axis has size 0")
	}

	return &RuleTargetAxis{
		DeviceName:    device_name,
		Device:        device,
		Axis:          axis,
		Inverted:      inverted,
		DeadzoneStart: deadzoneStart,
		DeadzoneEnd:   deadzoneEnd,
		deadzoneSize:  deadzoneSize,
		axisSize:      axisSize,
	}, nil
}

// NormalizeValue takes a raw input value and converts it to a value suitable for output.
//
// Axis inputs are normalized to the full signed int32 range to match the virtual device's axis
// characteristics.
//
// Typically this function is called after RuleTargetAxis.MatchEvent, which checks whether we are
// in the deadzone, among other things.
func (target *RuleTargetAxis) NormalizeValue(value int32) int32 {
	axisStrength := target.GetAxisStrength(value)
	return LerpInt(AxisValueMin, AxisValueMax, axisStrength)
}

func (target *RuleTargetAxis) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	value = Clamp(value, AxisValueMin, AxisValueMax)
	return &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  target.Axis,
		Value: value,
	}
}

func (target *RuleTargetAxis) MatchEvent(device RuleTargetDevice, event *evdev.InputEvent) bool {
	return target.MatchEventDeviceAndCode(device, event) &&
		!target.InDeadZone(event.Value)
}

// TODO: Add tests
func (target *RuleTargetAxis) MatchEventDeviceAndCode(device RuleTargetDevice, event *evdev.InputEvent) bool {
	return device == target.Device &&
		event.Type == evdev.EV_ABS &&
		event.Code == target.Axis
}

// TODO: Add tests
func (target *RuleTargetAxis) InDeadZone(value int32) bool {
	return value >= target.DeadzoneStart && value <= target.DeadzoneEnd
}

// GetAxisStrength returns a float between 0.0 and 1.0, representing the proportional
// position along the axis' full range. (after factoring in deadzones)
// Calling this function with `value` inside the deadzone range will produce undefined behavior
func (target *RuleTargetAxis) GetAxisStrength(value int32) float64 {
	if value > target.DeadzoneEnd {
		value -= target.deadzoneSize
	}
	strength := float64(value) / float64(target.axisSize)
	if target.Inverted {
		strength = 1.0 - strength
	}
	return strength
}
