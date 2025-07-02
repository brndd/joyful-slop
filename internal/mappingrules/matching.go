package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

// eventFromTarget creates an outputtable event from a RuleTarget
func eventFromTarget(output RuleTarget, value int32) *evdev.InputEvent {
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
			// TODO: how would we invert axes?
		default:
			logger.Logf("Inverted rule for unknown event type '%d'. Not inverting value\n", event.Type)
		}
	}

	return value
}

func (rule SimpleMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent) *evdev.InputEvent {
	if device != rule.Input.Device ||
		event.Code != rule.Input.Code {
		return nil
	}

	return eventFromTarget(rule.Output, valueFromTarget(rule.Input, event))
}

func (rule ComboMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent) *evdev.InputEvent {
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
		rule.State--
	}
	if inputValue == 1 {
		rule.State++
	}
	targetState := len(rule.Inputs)
	if oldState == targetState-1 && rule.State == targetState {
		return eventFromTarget(rule.Output, 1)
	}
	if oldState == targetState && rule.State == targetState-1 {
		return eventFromTarget(rule.Output, 0)
	}
	return nil
}
