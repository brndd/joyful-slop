package mappingrules

import (
	"testing"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/suite"
)

type MappingRuleAxisToButtonTests struct {
	suite.Suite
	inputDevice  *InputDeviceMock
	inputRule    *RuleTargetAxis
	outputDevice *evdev.InputDevice
	outputRule   *RuleTargetButton
	mode         *string
	base         MappingRuleBase
}

func TestRunnerMappingRuleAxisToButtonTests(t *testing.T) {
	suite.Run(t, new(MappingRuleAxisToButtonTests))
}

// buildTimerRule creates a MappingRuleAxisToButton with a mocked clock
func (t *MappingRuleAxisToButtonTests) buildTimerRule(
	repeatMin,
	repeatMax int,
	nextEvent time.Duration) (*MappingRuleAxisToButton, *clockwork.FakeClock) {

	mockClock := clockwork.NewFakeClock()
	testRule := t.buildRule(repeatMin, repeatMax)
	testRule.clock = mockClock
	testRule.lastEvent = testRule.clock.Now()
	testRule.nextEvent = nextEvent
	if nextEvent != NoNextEvent {
		testRule.active = true
	}
	return testRule, mockClock
}

// Todo: don't love this repeated logic...
func (t *MappingRuleAxisToButtonTests) buildRule(repeatMin, repeatMax int) *MappingRuleAxisToButton {
	return &MappingRuleAxisToButton{
		MappingRuleBase: t.base,
		Input:           t.inputRule,
		Output:          t.outputRule,
		RepeatRateMin:   repeatMin,
		RepeatRateMax:   repeatMax,
		lastEvent:       time.Now(),
		nextEvent:       NoNextEvent,
		repeat:          repeatMin != 0 && repeatMax != 0,
		pressed:         false,
		active:          false,
		clock:           clockwork.NewRealClock(),
	}
}

func (t *MappingRuleAxisToButtonTests) SetupTest() {
	mode := "*"
	t.mode = &mode
	t.inputDevice = new(InputDeviceMock)
	t.inputDevice.On("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
		evdev.ABS_X: {
			Minimum: 0,
			Maximum: 10000,
		},
	}, nil)
	t.inputRule, _ = NewRuleTargetAxis("test-input", t.inputDevice, evdev.ABS_X, false, []Deadzone{{Start: 0, End: 1000}})

	t.outputDevice = &evdev.InputDevice{}
	t.outputRule, _ = NewRuleTargetButton("test-output", t.outputDevice, evdev.ABS_X, false)
	t.base = NewMappingRuleBase("", []string{"*"})
}

func (t *MappingRuleAxisToButtonTests) TestMatchEvent() {

	// A valid input should set a nextevent
	t.Run("No Repeat", func() {
		testRule := t.buildRule(0, 0)

		t.Run("Valid Input", func() {
			testRule.MatchEvent(t.inputDevice, &evdev.InputEvent{
				Type:  evdev.EV_ABS,
				Code:  evdev.ABS_X,
				Value: 1001,
			}, t.mode)
			t.NotEqual(NoNextEvent, testRule.nextEvent)
		})

		t.Run("Deadzone Input", func() {
			testRule.MatchEvent(t.inputDevice, &evdev.InputEvent{
				Type:  evdev.EV_ABS,
				Code:  evdev.ABS_X,
				Value: 500,
			}, t.mode)
			t.Equal(NoNextEvent, testRule.nextEvent)
		})
	})

	t.Run("Repeat", func() {
		testRule := t.buildRule(750, 250)
		testRule.MatchEvent(t.inputDevice, &evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_X,
			Value: 10000,
		}, t.mode)
		t.Equal(time.Duration(250*time.Millisecond), testRule.nextEvent)

		testRule.MatchEvent(t.inputDevice, &evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_X,
			Value: 1001,
		}, t.mode)
		// Allow leeway since time passes during the test
		t.True(testRule.nextEvent > time.Duration(650*time.Millisecond))

		testRule.MatchEvent(t.inputDevice, &evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_X,
			Value: 5500,
		}, t.mode)
		// Allow up to 50 ms leeway since time passes during the test
		t.InDelta(time.Duration(500*time.Millisecond), testRule.nextEvent, 50000000)
	})
}

func (t *MappingRuleAxisToButtonTests) TestTimerEvent() {
	t.Run("No Repeat", func() {
		// Get event if called immediately
		t.Run("Event is available immediately", func() {
			testRule, _ := t.buildTimerRule(0, 0, 0)

			event := testRule.TimerEvent()

			t.EqualValues(1, event.Value)
			t.Equal(true, testRule.pressed)
		})

		// Off event on second call
		t.Run("Event emits off on second call", func() {
			testRule, _ := t.buildTimerRule(0, 0, 0)

			testRule.TimerEvent()
			event := testRule.TimerEvent()

			t.EqualValues(0, event.Value)
			t.Equal(false, testRule.pressed)
		})

		// No further event, even if we wait a while
		t.Run("Additional events are not emitted while still active.", func() {
			testRule, mockClock := t.buildTimerRule(0, 0, 0)

			testRule.TimerEvent()
			testRule.TimerEvent()

			mockClock.Advance(10 * time.Millisecond)
			event := testRule.TimerEvent()
			t.Nil(event)
			t.Equal(false, testRule.pressed)
		})
	})

	t.Run("Repeat", func() {
		t.Run("No event if called immediately", func() {
			testRule, _ := t.buildTimerRule(100, 10, 50*time.Millisecond)
			event := testRule.TimerEvent()
			t.Nil(event)
		})

		t.Run("No event after 49ms", func() {
			testRule, mockClock := t.buildTimerRule(100, 10, 50*time.Millisecond)
			mockClock.Advance(49 * time.Millisecond)

			event := testRule.TimerEvent()

			t.Nil(event)
		})

		t.Run("Event after 50ms", func() {
			testRule, mockClock := t.buildTimerRule(100, 10, 50*time.Millisecond)
			mockClock.Advance(50 * time.Millisecond)

			event := testRule.TimerEvent()

			t.EqualValues(1, event.Value)
			t.Equal(true, testRule.pressed)
		})

		t.Run("Additional event at 100ms", func() {
			testRule, mockClock := t.buildTimerRule(100, 10, 50*time.Millisecond)

			mockClock.Advance(50 * time.Millisecond)
			testRule.TimerEvent()
			testRule.TimerEvent()

			mockClock.Advance(50 * time.Millisecond)
			event := testRule.TimerEvent()

			t.NotNil(event)
		})
	})
}
