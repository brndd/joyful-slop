package mappingrules

import "github.com/holoplot/go-evdev"

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleAxis struct {
	MappingRuleBase
	Input  *RuleTargetAxis
	Output *RuleTargetAxis
}

func NewMappingRuleAxis(base MappingRuleBase, input *RuleTargetAxis, output *RuleTargetAxis) *MappingRuleAxis {
	return &MappingRuleAxis{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}
}

func (rule *MappingRuleAxis) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEvent(device, event) {
		return nil, nil
	}

	return rule.Output.Device, rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
