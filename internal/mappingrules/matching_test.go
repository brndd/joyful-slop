package mappingrules

import (
	"testing"

	"github.com/holoplot/go-evdev"
)

func TestSimpleRuleMatchEvent(t *testing.T) {
	inputDevice := &evdev.InputDevice{}
	wrongInputDevice := &evdev.InputDevice{}
	outputDevice := &evdev.InputDevice{}

	rule := &SimpleMappingRule{
		MappingRuleBase: MappingRuleBase{
			Output: &RuleTargetButton{
				RuleTargetBase{
					DeviceName: "test_output",
					Device:     outputDevice,
					Code:       evdev.BTN_TRIGGER,
				},
			},
			Modes: []string{"*"},
		},
		Input: &RuleTargetButton{
			RuleTargetBase{
				DeviceName: "test_input",
				Device:     inputDevice,
				Code:       evdev.BTN_TRIGGER,
			},
		},
	}
	mode := "*"

	// A matching input event should produce an output event
	event := rule.MatchEvent(inputDevice, &evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, &mode)
	outputEvent := &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  evdev.BTN_TRIGGER,
		Value: 1,
	}

	if event == nil || *event != *outputEvent {
		t.Errorf("Expected event to match %v, but got %v", outputEvent, event)
	}

	// if event.Type != outputEvent.Type ||
	// 	event.Code != outputEvent.Code ||
	// 	event.Value != outputEvent.Value {
	// 	t.Errorf("Expected event to match %v, but got %v", outputEvent, event)
	// }

	// An input event from the wrong device should produce a nil event
	event = rule.MatchEvent(wrongInputDevice, &evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, &mode)
	if event != nil {
		t.Errorf("Expected event not to match, but got non-nil event %v", event)
	}

	// An input event from the wrong device should produce a nil event
	event = rule.MatchEvent(wrongInputDevice, &evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, &mode)
	if event != nil {
		t.Errorf("Expected event not to match, but got non-nil event %v", event)
	}

	// TODO: test inversion, and everything else...
}
