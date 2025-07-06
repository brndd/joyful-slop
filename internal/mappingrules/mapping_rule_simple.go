package mappingrules

import "github.com/holoplot/go-evdev"

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleSimple struct {
	MappingRuleBase
	Input RuleTarget
}

func (rule *MappingRuleSimple) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	if device != rule.Input.GetDevice() ||
		event.Code != rule.Input.GetCode() {
		return nil
	}

	return rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
