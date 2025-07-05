package mappingrules

import (
	"slices"
	"time"

	"github.com/holoplot/go-evdev"
)

type MappingRuleBase struct {
	Name   string
	Output RuleTarget
	Modes  []string
}

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleSimple struct {
	MappingRuleBase
	Input RuleTarget
}

// A Combo Mapping Rule can require multiple physical button presses for a single output button
type MappingRuleCombo struct {
	MappingRuleBase
	Inputs []RuleTarget
	State  int
}

type MappingRuleLatched struct {
	MappingRuleBase
	Input RuleTarget
	State bool
}

// TODO: How are we going to implement this? It needs to operate on a timer...
type MappingRuleProportionalAxis struct {
	MappingRuleBase
	Input       *RuleTargetAxis
	Output      *RuleTargetButton
	Sensitivity int32
	LastValue   int32
	LastEvent   time.Time
}

func (rule *MappingRuleBase) OutputName() string {
	return rule.Output.GetDeviceName()
}

func (rule *MappingRuleBase) modeCheck(mode *string) bool {
	if rule.Modes[0] == "*" {
		return true
	}
	return slices.Contains(rule.Modes, *mode)
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

func (rule *MappingRuleLatched) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	if device != rule.Input.GetDevice() ||
		event.Code != rule.Input.GetCode() ||
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

func (rule *MappingRuleProportionalAxis) MatchEvent(device *evdev.InputDevice, event *evdev.InputEvent, mode *string) *evdev.InputEvent {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil
	}

	if device != rule.Input.GetDevice() ||
		event.Code != rule.Input.GetCode() {
		return nil
	}

	// set the last value to the normalized input value
	rule.LastValue = rule.Input.NormalizeValue(event.Value)
	return nil
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
