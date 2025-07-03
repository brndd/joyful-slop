package mappingrules

import (
	"slices"

	"git.annabunches.net/annabunches/joyful/internal/logger"
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

// eventFromTarget creates an outputtable event from a RuleTarget
func eventFromTarget(output RuleTarget, value int32, mode *string) *evdev.InputEvent {
	if len(output.ModeSelect) > 0 {
		index := 0
		if currentMode := slices.Index(output.ModeSelect, *mode); currentMode != -1 {
			// find the next mode
			index = (currentMode + 1) % len(output.ModeSelect)
		}

		*mode = output.ModeSelect[index]
		return nil
	}
	return &evdev.InputEvent{
		Type:  output.Type,
		Code:  output.Code,
		Value: value,
	}
}

// valueFromTarget determines the value to output from an input specification,given a RuleTarget's constraints
func valueFromTarget(rule RuleTarget, event *evdev.InputEvent) int32 {
	// how we process inverted rules depends on the event type
	value := event.Value
	if rule.Inverted {
		switch rule.Type {
		case evdev.EV_KEY:
			if value == 0 {
				value = 1
			} else {
				value = 0
			}
		case evdev.EV_ABS:
			logger.Logf("STUB: Inverting axes is not yet implemented.")
		default:
			logger.Logf("Inverted rule for unknown event type '%d'. Not inverting value", event.Type)
		}
	}

	return value
}

func (rule *SimpleMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	if event.Type != evdev.EV_ABS {
		logger.Logf("DEBUG: mode check passed for rule '%s'. Mode '%s' modes '%v'", rule.Name, *mode, rule.Modes)
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
		valueFromTarget(rule.Input, event) == 0 {
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

	return eventFromTarget(rule.Output, value, mode)
}
