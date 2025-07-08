package mappingrules

import "github.com/holoplot/go-evdev"

type MappingRuleModeSelect struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetModeSelect
}

func NewMappingRuleModeSelect(
	base MappingRuleBase,
	input *RuleTargetButton,
	output *RuleTargetModeSelect,
) *MappingRuleModeSelect {

	return &MappingRuleModeSelect{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}
}

func (rule *MappingRuleModeSelect) MatchEvent(
	device *evdev.InputDevice,
	event *evdev.InputEvent,
	mode *string) (*evdev.InputDevice, *evdev.InputEvent) {

	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Button ||
		rule.Input.NormalizeValue(event.Value) == 0 {
		return nil, nil
	}

	return nil, rule.Output.CreateEvent(event.Value, mode)
}
