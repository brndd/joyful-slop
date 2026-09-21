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

func TestButtonTempoTapEmitsCompleteMultiButtonPulse(t *testing.T) {
	rule, clock, inputDevice, mode := newButtonTempoTestRule(t)
	require.Empty(t, rule.MatchEvents(inputDevice, buttonInput(1), mode))
	clock.Advance(499 * time.Millisecond)
	require.Empty(t, rule.TimerEvents(mode))

	events := rule.MatchEvents(inputDevice, buttonInput(0), mode)
	require.Len(t, events, 4)
	require.Equal(t, []int32{1, 1, 0, 0}, []int32{events[0].Event.Value, events[1].Event.Value, events[2].Event.Value, events[3].Event.Value})
	require.Equal(t, "tap-mode", *mode)
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
	require.Len(t, events, 4)
	require.Equal(t, []int32{1, 1, 0, 0}, []int32{events[0].Event.Value, events[1].Event.Value, events[2].Event.Value, events[3].Event.Value})
	require.Equal(t, "hold-mode", *mode)
	require.Empty(t, rule.TimerEvents(mode))
}
