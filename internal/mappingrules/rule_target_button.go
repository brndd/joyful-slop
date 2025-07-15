package mappingrules

import "github.com/holoplot/go-evdev"

type RuleTargetButton struct {
	DeviceName string
	Device     *evdev.InputDevice
	Button     evdev.EvCode
	Inverted   bool
}

func NewRuleTargetButton(device_name string, device *evdev.InputDevice, code evdev.EvCode, inverted bool) (*RuleTargetButton, error) {
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

func (target *RuleTargetButton) MatchEvent(device RuleTargetDevice, event *evdev.InputEvent) bool {
	return device == target.Device &&
		event.Type == evdev.EV_KEY &&
		event.Code == target.Button
}
