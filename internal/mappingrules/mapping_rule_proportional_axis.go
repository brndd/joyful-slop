package mappingrules

import (
	"time"

	"github.com/holoplot/go-evdev"
)

// TODO: How are we going to implement this? It needs to operate on a timer...
type MappingRuleProportionalAxis struct {
	MappingRuleBase
	Input       *RuleTargetAxis
	Output      *RuleTargetButton
	Sensitivity int32
	LastValue   int32
	LastEvent   time.Time
}

func (rule *MappingRuleProportionalAxis) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	if device != rule.Input.GetDevice() ||
		event.Code != rule.Input.GetCode() {
		return nil, nil
	}

	// set the last value to the normalized input value
	rule.LastValue = rule.Input.NormalizeValue(event.Value)
	return nil, nil
}

// TimerEvent returns an event when enough time has passed (compared to the last recorded axis value)
// to emit an event.
func (rule *MappingRuleProportionalAxis) TimerEvent() *evdev.InputEvent {
	// This is tighter coupling than we'd like, but it will do for now.
	// TODO: maybe it would be better to just be more declarative about event types and their inputs and outputs.
	if rule.LastValue < rule.Input.AxisStart {
		rule.LastEvent = time.Now()
		return nil
	}

	// calculate target time until next event press
	//	nextEvent := rule.LastEvent + (rule.LastValue)

	// TODO: figure out what the condition should be
	if false {
		// TODO: emit event
		rule.LastEvent = time.Now()
	}

	return nil
}
