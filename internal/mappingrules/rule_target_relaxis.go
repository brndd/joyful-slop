package mappingrules

import (
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/eventcodes"
	"github.com/holoplot/go-evdev"
)

type RuleTargetRelaxis struct {
	DeviceName string
	Device     Device
	Axis       evdev.EvCode
}

func NewRuleTargetRelaxisFromConfig(targetConfig configparser.RuleTargetConfigRelaxis, devs map[string]Device) (*RuleTargetRelaxis, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	eventCode, err := eventcodes.ParseCode(targetConfig.Axis, eventcodes.CodePrefixRelaxis)
	if err != nil {
		return nil, err
	}

	return NewRuleTargetRelaxis(
		targetConfig.Device,
		device,
		eventCode,
	)
}

func NewRuleTargetRelaxis(deviceName string,
	device Device,
	axis evdev.EvCode) (*RuleTargetRelaxis, error) {

	return &RuleTargetRelaxis{
		DeviceName: deviceName,
		Device:     device,
		Axis:       axis,
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
