package mappingrules

import (
	"testing"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/require"
)

func newButtonTempoTestRule(t *testing.T) (*MappingRuleButtonTempo, *clockwork.FakeClock, *InputDeviceMock, *string) {
	t.Helper()
	inputDevice := new(InputDeviceMock)
	outputDevice := &evdev.InputDevice{}
	input, err := NewRuleTargetButton("input", inputDevice, evdev.BTN_TRIGGER, false)
	require.NoError(t, err)
	first, err := NewRuleTargetButton("output", outputDevice, evdev.BTN_THUMB, false)
	require.NoError(t, err)
	second, err := NewRuleTargetButton("output", outputDevice, evdev.BTN_THUMB2, false)
	require.NoError(t, err)
	clock := clockwork.NewFakeClock()
	rule := &MappingRuleButtonTempo{
		MappingRuleBase: NewMappingRuleBase("", nil),
		Input:           input,
		Threshold:       500 * time.Millisecond,
		tap:             buttonTempoBranch{outputs: []*RuleTargetButton{first, second}, mode: "tap-mode"},
		hold:            buttonTempoBranch{outputs: []*RuleTargetButton{first, second}, mode: "hold-mode"},
		clock:           clock,
	}
	mode := "base"
	return rule, clock, inputDevice, &mode
}

func buttonInput(value int32) *evdev.InputEvent {
	return &evdev.InputEvent{Type: evdev.EV_KEY, Code: evdev.BTN_TRIGGER, Value: value}
}

func TestButtonTempoTapSeparatesPressAndRelease(t *testing.T) {
	rule, clock, inputDevice, mode := newButtonTempoTestRule(t)
	require.Empty(t, rule.MatchEvents(inputDevice, buttonInput(1), mode))
	clock.Advance(499 * time.Millisecond)
	require.Empty(t, rule.TimerEvents(mode))

	events := rule.MatchEvents(inputDevice, buttonInput(0), mode)
	require.Len(t, events, 2)
	require.Equal(t, []int32{1, 1}, []int32{events[0].Event.Value, events[1].Event.Value})
	require.Equal(t, "tap-mode", *mode)

	clock.Advance(buttonTempoPulseDuration - time.Millisecond)
	require.Empty(t, rule.TimerEvents(mode))
	clock.Advance(time.Millisecond)
	events = rule.TimerEvents(mode)
	require.Len(t, events, 2)
	require.Equal(t, []int32{0, 0}, []int32{events[0].Event.Value, events[1].Event.Value})
}

func TestButtonTempoHoldFiresOnceAndReleasesOnInputRelease(t *testing.T) {
	rule, clock, inputDevice, mode := newButtonTempoTestRule(t)
	rule.MatchEvents(inputDevice, buttonInput(1), mode)
	clock.Advance(500 * time.Millisecond)

	events := rule.TimerEvents(mode)
	require.Len(t, events, 2)
	require.EqualValues(t, 1, events[0].Event.Value)
	require.Equal(t, "hold-mode", *mode)
	require.Empty(t, rule.TimerEvents(mode), "hold action must only fire once")

	events = rule.MatchEvents(inputDevice, buttonInput(0), mode)
	require.Len(t, events, 2)
	require.EqualValues(t, 0, events[0].Event.Value)
	require.EqualValues(t, 0, events[1].Event.Value)
	require.Equal(t, "hold-mode", *mode)
}

func TestButtonTempoReleaseAtThresholdIsHold(t *testing.T) {
	rule, clock, inputDevice, mode := newButtonTempoTestRule(t)
	rule.MatchEvents(inputDevice, buttonInput(1), mode)
	clock.Advance(500 * time.Millisecond)

	events := rule.MatchEvents(inputDevice, buttonInput(0), mode)
	require.Len(t, events, 2)
	require.Equal(t, []int32{1, 1}, []int32{events[0].Event.Value, events[1].Event.Value})
	require.Equal(t, "hold-mode", *mode)

	clock.Advance(buttonTempoPulseDuration)
	events = rule.TimerEvents(mode)
	require.Len(t, events, 2)
	require.Equal(t, []int32{0, 0}, []int32{events[0].Event.Value, events[1].Event.Value})
}

func TestButtonTempoNewPressDrainsPendingRelease(t *testing.T) {
	rule, clock, inputDevice, mode := newButtonTempoTestRule(t)
	rule.MatchEvents(inputDevice, buttonInput(1), mode)
	clock.Advance(100 * time.Millisecond)
	events := rule.MatchEvents(inputDevice, buttonInput(0), mode)
	require.Len(t, events, 2)

	events = rule.MatchEvents(inputDevice, buttonInput(1), mode)
	require.Len(t, events, 2)
	require.Equal(t, []int32{0, 0}, []int32{events[0].Event.Value, events[1].Event.Value})
}

func TestButtonTempoExternalModeChangeReleasesHeldOutputs(t *testing.T) {
	rule, clock, inputDevice, mode := newButtonTempoTestRule(t)
	rule.MappingRuleBase = NewMappingRuleBase("", []string{"base"})
	rule.MatchEvents(inputDevice, buttonInput(1), mode)
	clock.Advance(500 * time.Millisecond)

	presses := rule.TimerEvents(mode)
	require.Len(t, presses, 2)
	require.Empty(t, rule.ModeChanged(*mode, true), "the initiating transition must preserve the hold")

	releases := rule.ModeChanged("other", false)
	require.Len(t, releases, 2)
	require.Equal(t, []int32{0, 0}, []int32{releases[0].Event.Value, releases[1].Event.Value})
	require.False(t, rule.active)
	require.False(t, rule.held)
}
