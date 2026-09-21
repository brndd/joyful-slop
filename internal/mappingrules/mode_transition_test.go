package mappingrules

import (
	"testing"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/require"
)

func TestButtonModeChangeReleasesActiveOutput(t *testing.T) {
	inputDevice := &evdev.InputDevice{}
	outputDevice := &evdev.InputDevice{}
	input, _ := NewRuleTargetButton("", inputDevice, evdev.BTN_TRIGGER, false)
	output, _ := NewRuleTargetButton("", outputDevice, evdev.BTN_TOP, false)
	rule := &MappingRuleButton{MappingRuleBase: NewMappingRuleBase("", []string{"a"}), Input: input, Output: output}
	mode := "a"

	_, press := rule.MatchEvent(inputDevice, &evdev.InputEvent{Type: evdev.EV_KEY, Code: evdev.BTN_TRIGGER, Value: 1}, &mode)
	require.EqualValues(t, 1, press.Value)

	events := rule.ModeChanged("b", false)
	require.Len(t, events, 1)
	require.Same(t, outputDevice, events[0].Device)
	require.EqualValues(t, 0, events[0].Event.Value)
	require.Empty(t, rule.ModeChanged("b", false))
}

func TestHatModeChangeCentersActiveOutput(t *testing.T) {
	inputDevice := &evdev.InputDevice{}
	outputDevice := &evdev.InputDevice{}
	rule := &MappingRuleHat{
		MappingRuleBase: NewMappingRuleBase("", []string{"a"}),
		Input:           &RuleTargetHat{Device: inputDevice, Hat: evdev.ABS_HAT0X},
		Output:          &RuleTargetHat{Device: outputDevice, Hat: evdev.ABS_HAT1X},
	}
	mode := "a"

	_, direction := rule.MatchEvent(inputDevice, &evdev.InputEvent{Type: evdev.EV_ABS, Code: evdev.ABS_HAT0X, Value: -1}, &mode)
	require.EqualValues(t, -1, direction.Value)

	events := rule.ModeChanged("b", false)
	require.Len(t, events, 1)
	require.EqualValues(t, 0, events[0].Event.Value)
}

func TestAxisToButtonModeChangeStopsAndReleases(t *testing.T) {
	inputDevice := NewInputDeviceMock()
	inputDevice.Stub("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
		evdev.ABS_X: {Minimum: 0, Maximum: 100},
	}, nil)
	input, err := NewRuleTargetAxis("", inputDevice, evdev.ABS_X, false, []Deadzone{{Start: 0, End: 10}})
	require.NoError(t, err)
	outputDevice := &evdev.InputDevice{}
	output, _ := NewRuleTargetButton("", outputDevice, evdev.BTN_TRIGGER, false)
	rule := &MappingRuleAxisToButton{
		MappingRuleBase: NewMappingRuleBase("", []string{"a"}),
		Input:           input,
		Output:          output,
		Hold:            true,
		nextEvent:       NoNextEvent,
	}
	mode := "a"

	_, press := rule.MatchEvent(inputDevice, &evdev.InputEvent{Type: evdev.EV_ABS, Code: evdev.ABS_X, Value: 100}, &mode)
	require.EqualValues(t, 1, press.Value)

	events := rule.ModeChanged("b", false)
	require.Len(t, events, 1)
	require.EqualValues(t, 0, events[0].Event.Value)
	require.False(t, rule.pressed)
	require.Equal(t, NoNextEvent, rule.nextEvent)
}

func TestModeChangeCancelsRepeatingRules(t *testing.T) {
	buttonRule := &MappingRuleAxisToButton{
		MappingRuleBase: NewMappingRuleBase("", []string{"a"}),
		nextEvent:       time.Millisecond,
		active:          true,
	}
	relaxisRule := &MappingRuleAxisToRelaxis{
		MappingRuleBase: NewMappingRuleBase("", []string{"a"}),
		nextEvent:       time.Millisecond,
	}

	require.Empty(t, buttonRule.ModeChanged("b", false))
	require.Equal(t, NoNextEvent, buttonRule.nextEvent)
	require.False(t, buttonRule.active)
	require.Empty(t, relaxisRule.ModeChanged("b", false))
	require.Equal(t, NoNextEvent, relaxisRule.nextEvent)
}

func TestStatefulRulesResetTheirOutputs(t *testing.T) {
	outputDevice := &evdev.InputDevice{}
	output, _ := NewRuleTargetButton("", outputDevice, evdev.BTN_TRIGGER, false)
	combo := &MappingRuleButtonCombo{Inputs: []*RuleTargetButton{{}, {}}, Output: output, State: 2}
	latched := &MappingRuleButtonLatched{Output: output, State: true}
	mode := "a"

	comboEvents := combo.Reset(&mode)
	require.Len(t, comboEvents, 1)
	require.Equal(t, 0, combo.State)
	require.EqualValues(t, 0, comboEvents[0].Event.Value)

	latchedEvents := latched.Reset(&mode)
	require.Len(t, latchedEvents, 1)
	require.False(t, latched.State)
	require.EqualValues(t, 0, latchedEvents[0].Event.Value)
}
