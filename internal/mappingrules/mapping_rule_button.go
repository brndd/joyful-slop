package mappingrules

import "github.com/holoplot/go-evdev"

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleButton struct {
	MappingRuleBase
	Input  *RuleTargetButton
	Output *RuleTargetButton
}

func (rule *MappingRuleButton) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.GetDevice() ||
		event.Code != rule.Input.GetCode() {
		return nil, nil
	}

	return rule.Output.Device, rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
