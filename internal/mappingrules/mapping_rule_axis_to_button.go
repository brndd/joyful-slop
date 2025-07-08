package mappingrules

import (
	"time"

	"github.com/holoplot/go-evdev"
)

// TODO: This whole file is still WIP
type MappingRuleAxisToButton struct {
	MappingRuleBase
	Input          *RuleTargetAxis
	Output         *RuleTargetButton
	RepeatSpeedMin int32
	RepeatSpeedMax int32
	lastValue      int32
	lastEvent      time.Time
}

func (rule *MappingRuleAxisToButton) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEvent(device, event) {
		return nil, nil
	}

	// set the last value to the normalized input value
	rule.lastValue = rule.Input.NormalizeValue(event.Value)
	return nil, nil
}

// TimerEvent returns an event when enough time has passed (compared to the last recorded axis value)
// to emit an event.
func (rule *MappingRuleAxisToButton) TimerEvent() *evdev.InputEvent {
	// This is tighter coupling than we'd like, but it will do for now.
	// TODO: maybe it would be better to just be more declarative about event types and their inputs and outputs.
	if rule.lastValue < rule.Input.DeadzoneStart {
		rule.lastEvent = time.Now()
		return nil
	}

	// calculate target time until next event press
	//	nextEvent := rule.LastEvent + (rule.LastValue)

	// TODO: figure out what the condition should be
	if false {
		// TODO: emit event
		rule.lastEvent = time.Now()
	}

	return nil
}
