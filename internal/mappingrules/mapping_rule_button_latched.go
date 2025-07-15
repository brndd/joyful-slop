package mappingrules

import "github.com/holoplot/go-evdev"

type MappingRuleButtonLatched struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetButton
	State  bool
}

func NewMappingRuleButtonLatched(
	base MappingRuleBase,
	input *RuleTargetButton,
	output *RuleTargetButton) *MappingRuleButtonLatched {

	return &MappingRuleButtonLatched{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
		State:           false,
	}
}

func (rule *MappingRuleButtonLatched) MatchEvent(device RuleTargetDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Button ||
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
