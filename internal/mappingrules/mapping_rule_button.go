package mappingrules

import "github.com/holoplot/go-evdev"

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleButton struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetButton
}

func NewMappingRuleButton(
	base MappingRuleBase,
	input *RuleTargetButton,
	output *RuleTargetButton) *MappingRuleButton {

	return &MappingRuleButton{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}
}

func (rule *MappingRuleButton) MatchEvent(device RuleTargetDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Button {
		return nil, nil
	}

	return rule.Output.Device, rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
