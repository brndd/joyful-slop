package mappingrules

import (
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/jonboulle/clockwork"
)

// MappingRuleAxisToButton represents a rule that converts an axis input into a (potentially repeating)
// button output.
type MappingRuleAxisToButton struct {
	MappingRuleBase
	Input         *RuleTargetAxis
	Output        *RuleTargetButton
	RepeatRateMin int
	RepeatRateMax int
	nextEvent     time.Duration
	lastEvent     time.Time
	repeat        bool
	pressed       bool // "pressed" indicates that we've sent the output button signal, but still need to send the button release
	active        bool // "active" is true whenever the input is not in a deadzone
	clock         clockwork.Clock
}

func NewMappingRuleAxisToButton(base MappingRuleBase, input *RuleTargetAxis, output *RuleTargetButton, repeatRateMin, repeatRateMax int) *MappingRuleAxisToButton {
	return &MappingRuleAxisToButton{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
		RepeatRateMin:   repeatRateMin,
		RepeatRateMax:   repeatRateMax,
		lastEvent:       time.Now(),
		nextEvent:       NoNextEvent,
		repeat:          repeatRateMin != 0 && repeatRateMax != 0,
		pressed:         false,
		active:          false,
		clock:           clockwork.NewRealClock(),
	}
}

func (rule *MappingRuleAxisToButton) MatchEvent(device RuleTargetDevice, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {

	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEventDeviceAndCode(device, event) {
		return nil, nil
	}

	// If we're inside the deadzone, unset the next event
	if rule.Input.InDeadZone(event.Value) {
		rule.nextEvent = NoNextEvent
		rule.active = false
		return nil, nil
	}

	// If we aren't repeating, we trigger the event immediately
	// We also only set this if active == false, so that only one
	// event can be emitted per "active" period
	if !rule.repeat && !rule.active {
		rule.nextEvent = 0
		rule.active = true
		return nil, nil
	}

	// use the axis value and the repeat rate to set a target time until the next event
	strength := 1.0 - rule.Input.GetAxisStrength(event.Value)
	rate := int64(LerpInt(rule.RepeatRateMax, rule.RepeatRateMin, strength))
	rule.nextEvent = time.Duration(rate * int64(time.Millisecond))
	rule.active = true

	return nil, nil
}

// TimerEvent returns an event when enough time has passed (compared to the last recorded axis value)
// to emit an event.
func (rule *MappingRuleAxisToButton) TimerEvent() *evdev.InputEvent {
	// If we pressed the button last tick, release it before doing anything else
	if rule.pressed {
		rule.pressed = false
		return rule.Output.CreateEvent(0, nil)
	}

	// If we should not emit another event,
	// we just update lastEvent for station keeping
	if rule.nextEvent == NoNextEvent {
		rule.lastEvent = rule.clock.Now()
		return nil
	}

	if rule.clock.Now().Compare(rule.lastEvent.Add(rule.nextEvent)) > -1 {
		rule.lastEvent = rule.clock.Now()
		rule.pressed = true

		// The default case here is to leave nextEvent at whatever
		// it has been set to by MatchEvent. Since nextEvent is a delta,
		// this will naturally cause the repeat to happen
		if !rule.repeat {
			rule.nextEvent = NoNextEvent
		}

		return rule.Output.CreateEvent(1, nil)
	}

	return nil
}

func (rule *MappingRuleAxisToButton) GetOutputDevice() *evdev.InputDevice {
	return rule.Output.Device
}
