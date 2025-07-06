package mappingrules

import "github.com/holoplot/go-evdev"

type MappingRuleLatched struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetButton
	State  bool
}

func (rule *MappingRuleLatched) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Code ||
		rule.Input.NormalizeValue(event.Value) == 0 {
		return nil, nil
	}

	// Input is pressed, so toggle state and emit event
	var value int32
	rule.State = !rule.State
	if rule.State {
		value = 1
	} else {
		value = 0
	}

	return rule.Output.Device, rule.Output.CreateEvent(value, mode)
}
