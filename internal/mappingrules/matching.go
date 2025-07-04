package mappingrules

import (
	"slices"

	"github.com/holoplot/go-evdev"
)

func (rule *MappingRuleBase) OutputName() string {
	return rule.Output.DeviceName
}

func (rule *MappingRuleBase) modeCheck(mode *string) bool {
	if len(rule.Modes) == 1 && rule.Modes[0] == "*" {
		return true
	}
	return slices.Contains(rule.Modes, *mode)
}

func (rule *SimpleMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Code {
		return nil
	}

	return eventFromTarget(rule.Output, valueFromTarget(rule.Input, event), mode)
}

func (rule *ComboMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	// Check each of the inputs, and if we find a match, proceed
	var match *RuleTarget
	for _, input := range rule.Inputs {
		if device == input.Device &&
			event.Code == input.Code {
			match = &input
		}
	}

	if match == nil {
		return nil
	}

	// Get the value and add/subtract it from State
	inputValue := valueFromTarget(*match, event)
	oldState := rule.State
	if inputValue == 0 {
		rule.State = max(rule.State-1, 0)
	}
	if inputValue == 1 {
		rule.State++
	}
	targetState := len(rule.Inputs)

	if oldState == targetState-1 && rule.State == targetState {
		return eventFromTarget(rule.Output, 1, mode)
	}
	if oldState == targetState && rule.State == targetState-1 {
		return eventFromTarget(rule.Output, 0, mode)
	}
	return nil
}

func (rule *LatchedMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	if device != rule.Input.Device ||
		event.Code != rule.Input.Code ||
		rule.Input.NormalizeValue(event.Value) == 0 {
		return nil
	}

	// Input is pressed, so toggle state and emit event
	var value int32
	rule.State = !rule.State
	if rule.State {
		value = 1
	} else {
		value = 0
	}

	return rule.Output.CreateEvent(value, mode)
}

func (rule *ProportionalAxisMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	// STUB
	return nil
}

// TimerEvent returns an event when enough time has passed (compared to the last recorded axis value)
// to emit an event.
func (rule *ProportionalAxisMappingRule) TimerEvent() *evdev.InputEvent {
	// STUB
	return nil
}
