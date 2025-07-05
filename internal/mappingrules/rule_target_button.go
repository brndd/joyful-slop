package mappingrules

import "github.com/holoplot/go-evdev"

type RuleTargetButton struct {
	RuleTargetBase
}

func NewRuleTargetButton(device_name string, device *evdev.InputDevice, code evdev.EvCode, inverted bool) *RuleTargetButton {
	return &RuleTargetButton{
		RuleTargetBase: NewRuleTargetBase(device_name, device, code, inverted),
	}
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

func (target *RuleTargetButton) CreateEvent(value int32, mode *string) *evdev.InputEvent {
	return &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  target.Code,
		Value: value,
	}
}
