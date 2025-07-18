package mappingrules

import (
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/jonboulle/clockwork"
)

// TODO: add tests

// MappingRuleAxisToRelaxis represents a rule that converts an axis input into a (potentially repeating)
// relative axis output. This is most commonly used to generate mouse output events
type MappingRuleAxisToRelaxis struct {
	MappingRuleBase
	Input         *RuleTargetAxis
	Output        *RuleTargetRelaxis
	RepeatRateMin int
	RepeatRateMax int
	Increment     int32
	nextEvent     time.Duration
	lastEvent     time.Time
	clock         clockwork.Clock
}

func NewMappingRuleAxisToRelaxis(
	base MappingRuleBase,
	input *RuleTargetAxis,
	output *RuleTargetRelaxis,
	repeatRateMin, repeatRateMax, increment int) *MappingRuleAxisToRelaxis {

	return &MappingRuleAxisToRelaxis{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
		RepeatRateMin:   repeatRateMin,
		RepeatRateMax:   repeatRateMax,
		Increment:       int32(increment),
		lastEvent:       time.Now(),
		nextEvent:       NoNextEvent,
		clock:           clockwork.NewRealClock(),
	}
}

func (rule *MappingRuleAxisToRelaxis) MatchEvent(
	device Device,
	event *evdev.InputEvent,
	mode *string) (*evdev.InputDevice, *evdev.InputEvent) {

	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEventDeviceAndCode(device, event) {
		return nil, nil
	}

	// If we're inside the deadzone, unset the next event
	if rule.Input.InDeadZone(event.Value) {
		rule.nextEvent = NoNextEvent
		return nil, nil
	}

	// If we aren't repeating, we trigger the event immediately
	// TODO: this still needs the pressed parameter...
	if rule.RepeatRateMin == 0 || rule.RepeatRateMax == 0 {
		rule.nextEvent = time.Millisecond
		return nil, nil
	}

	// use the axis value and the repeat rate to set a target time until the next event
	strength := 1.0 - rule.Input.GetAxisStrength(event.Value)
	rate := int64(LerpInt(rule.RepeatRateMax, rule.RepeatRateMin, strength))
	rule.nextEvent = time.Duration(rate * int64(time.Millisecond))
	return nil, nil
}

// TimerEvent returns an event when enough time has passed (compared to the last recorded axis value)
// to emit an event.
func (rule *MappingRuleAxisToRelaxis) TimerEvent() *evdev.InputEvent {
	// This indicates that we should not emit another event
	if rule.nextEvent == NoNextEvent {
		rule.lastEvent = rule.clock.Now()
		return nil
	}

	if rule.clock.Now().Compare(rule.lastEvent.Add(rule.nextEvent)) > -1 {
		rule.lastEvent = rule.clock.Now()
		return rule.Output.CreateEvent(rule.Increment, nil)
	}

	return nil
}

func (rule *MappingRuleAxisToRelaxis) GetOutputDevice() *evdev.InputDevice {
	return rule.Output.Device.(*evdev.InputDevice)
}
