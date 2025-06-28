package rules

import (
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
)

type KeyMappingRule interface {
	MatchEvent(*evdev.InputDevice, *evdev.InputEvent)
}

type SimpleMappingRule struct {
	Input  RuleTarget
	Output RuleTarget
}

type ComboMappingRule struct {
	Input  []RuleTarget
	Output RuleTarget
}

func (rule *SimpleMappingRule) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent) *evdev.InputEvent {
	if event.Type != evdev.EV_KEY ||
		device != rule.Input.Device ||
		event.Code != rule.Input.Code {
		return nil
	}

	// how we process inverted rules depends on the event type
	value := event.Value
	if rule.Input.Inverted {
		switch rule.Input.Type {
		case evdev.EV_KEY:
			if value == 0 {
				value = 1
			} else {
				value = 0
			}
		default:
			logger.Log("Tried to handle inverted rule for unknown event type. Skipping rule.")
			return nil
		}
	}

	return &evdev.InputEvent{
		Type:  rule.Output.Type,
		Code:  rule.Output.Code,
		Value: value,
	}
}
