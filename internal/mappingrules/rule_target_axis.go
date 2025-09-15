package mappingrules

import (
	"errors"
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/eventcodes"
	"github.com/holoplot/go-evdev"
)

type RuleTargetAxis struct {
	DeviceName   string
	Device       Device
	Axis         evdev.EvCode
	Inverted     bool
	Deadzones    []Deadzone
	OutputMin    int32
	OutputMax    int32
	axisSize     int32
	deadzoneSize int32
}

func NewRuleTargetAxisFromConfig(targetConfig configparser.RuleTargetConfigAxis, devs map[string]Device) (*RuleTargetAxis, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	eventCode, err := eventcodes.ParseCode(targetConfig.Axis, eventcodes.CodePrefixAxis)
	if err != nil {
		return nil, err
	}

	deadzones := make([]Deadzone, 0)
	for _, dzConfig := range targetConfig.Deadzones {
		dz, err := NewDeadzoneFromConfig(dzConfig, device, eventCode)
		if err != nil {
			return nil, err
		}
		deadzones = append(deadzones, dz)
	}

	return NewRuleTargetAxis(
		targetConfig.Device,
		device,
		eventCode,
		targetConfig.Inverted,
		deadzones,
	)
}

func NewRuleTargetAxis(device_name string,
	device Device,
	axis evdev.EvCode,
	inverted bool,
	deadzones []Deadzone) (*RuleTargetAxis, error) {

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

	deadzoneSize := CalculateDeadzoneSize(deadzones)

	// Our output range is limited to 16 bits, but we represent values internally with 32 bits.
	// As a result, we shouldn't need to worry about integer overruns
	axisSize := info[axis].Maximum - info[axis].Minimum - deadzoneSize

	if axisSize == 0 {
		return nil, errors.New("axis has size 0")
	}

	return &RuleTargetAxis{
		DeviceName:   device_name,
		Device:       device,
		Axis:         axis,
		Inverted:     inverted,
		OutputMin:    AxisValueMin,
		OutputMax:    AxisValueMax,
		Deadzones:    deadzones,
		deadzoneSize: deadzoneSize,
		axisSize:     axisSize,
	}, nil
}

// NormalizeValue takes a raw input value and converts it to a value suitable for output.
//
// Axis inputs are normalized to the full signed int32 range to match the virtual device's axis
// characteristics.
//
// If the raw value is inside the deadzone, we either emit no event, or we emit the deadzoneValue.
// Typically this function is called after RuleTargetAxis.MatchEvent, which checks whether we are
// in the deadzone, among other things.
func (target *RuleTargetAxis) NormalizeValue(value int32) int32 {
	for _, dz := range target.Deadzones {
		state, dzValue := dz.Match(value)
		if state == DeadzoneEmit {
			return Clamp(dzValue, target.OutputMin, target.OutputMax)
		}
	}

	axisStrength := target.GetAxisStrength(value)
	return LerpInt(target.OutputMin, target.OutputMax, axisStrength)
}

func (target *RuleTargetAxis) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	value = Clamp(value, AxisValueMin, AxisValueMax)
	return &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  target.Axis,
		Value: value,
	}
}

func (target *RuleTargetAxis) MatchEvent(device Device, event *evdev.InputEvent) bool {
	return target.MatchEventDeviceAndCode(device, event) &&
		!target.InDeadZone(event.Value)
}

// TODO: Add tests
func (target *RuleTargetAxis) MatchEventDeviceAndCode(device Device, event *evdev.InputEvent) bool {
	return device == target.Device &&
		event.Type == evdev.EV_ABS &&
		event.Code == target.Axis
}

// InDeadZone checks each deadzone for whether the target value falls within it.
// If *any* non-emitting deadzone matches, we return true.
// TODO: Add tests
func (target *RuleTargetAxis) InDeadZone(value int32) bool {
	for _, dz := range target.Deadzones {
		state, _ := dz.Match(value)
		if state == DeadzoneNoEmit {
			return true
		}
	}
	return false
}

// GetAxisStrength returns a float between 0.0 and 1.0, representing the proportional
// position along the axis' full range. (after factoring in deadzones)
// Calling this function with `value` inside the deadzone range will produce undefined behavior
func (target *RuleTargetAxis) GetAxisStrength(value int32) float64 {
	adjValue := value
	for _, dz := range target.Deadzones {
		if value > dz.End {
			adjValue -= dz.Size
		}
	}
	strength := float64(adjValue) / float64(target.axisSize)
	if target.Inverted {
		strength = 1.0 - strength
	}
	return strength
}
