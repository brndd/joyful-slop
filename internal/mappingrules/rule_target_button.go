package mappingrules

import (
	"fmt"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/eventcodes"
	"github.com/holoplot/go-evdev"
)

type RuleTargetButton struct {
	DeviceName string
	Device     Device
	Button     evdev.EvCode
	Inverted   bool
}

func NewRuleTargetButtonFromConfig(targetConfig configparser.RuleTargetConfigButton, devs map[string]Device) (*RuleTargetButton, error) {
	device, ok := devs[targetConfig.Device]
	if !ok {
		return nil, fmt.Errorf("non-existent device '%s'", targetConfig.Device)
	}

	eventCode, err := eventcodes.ParseCodeButton(targetConfig.Button)
	if err != nil {
		return nil, err
	}

	return NewRuleTargetButton(
		targetConfig.Device,
		device,
		eventCode,
		targetConfig.Inverted,
	)
}

func NewRuleTargetButton(device_name string, device Device, code evdev.EvCode, inverted bool) (*RuleTargetButton, error) {
	return &RuleTargetButton{
		DeviceName: device_name,
		Device:     device,
		Button:     code,
		Inverted:   inverted,
	}, nil
}

func (target *RuleTargetButton) NormalizeValue(value int32) int32 {
	if target.Inverted {
		if value == 0 {
			return 1
		}
		return 0
	}
	return value
}

func (target *RuleTargetButton) CreateEvent(value int32, _ *string) *evdev.InputEvent {
	return &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  target.Button,
		Value: value,
	}
}

func (target *RuleTargetButton) MatchEvent(device Device, event *evdev.InputEvent) bool {
	return device == target.Device &&
		event.Type == evdev.EV_KEY &&
		event.Code == target.Button
}
