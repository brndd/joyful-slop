package mappingrules

import "github.com/holoplot/go-evdev"

// A Combo Mapping Rule can require multiple physical button presses for a single output button
type MappingRuleCombo struct {
	MappingRuleBase
	Inputs []RuleTarget
	State  int
}

func (rule *MappingRuleCombo) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	// Check each of the inputs, and if we find a match, proceed
	var match RuleTarget
	for _, input := range rule.Inputs {
		if device == input.GetDevice() &&
			event.Code == input.GetCode() {
			match = input
		}
	}

	if match == nil {
		return nil
	}

	// Get the value and add/subtract it from State
	inputValue := match.NormalizeValue(event.Value)
	oldState := rule.State
	if inputValue == 0 {
		rule.State = max(rule.State-1, 0)
	}
	if inputValue == 1 {
		rule.State++
	}
	targetState := len(rule.Inputs)

	if oldState == targetState-1 && rule.State == targetState {
		return rule.Output.CreateEvent(1, mode)
	}
	if oldState == targetState && rule.State == targetState-1 {
		return rule.Output.CreateEvent(0, mode)
	}
	return nil
}
