package mappingrules

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

type SimpleMappingRuleTests struct {
	suite.Suite
	inputDevice      *evdev.InputDevice
	wrongInputDevice *evdev.InputDevice
	outputDevice     *evdev.InputDevice
	mode             *string
	sampleRule       *SimpleMappingRule
	invertedRule     *SimpleMappingRule
}

func (t *SimpleMappingRuleTests) SetupTest() {
	t.inputDevice = &evdev.InputDevice{}
	t.wrongInputDevice = &evdev.InputDevice{}
	t.outputDevice = &evdev.InputDevice{}
	mode := "*"
	t.mode = &mode

	// TODO: implement a constructor function...
	t.sampleRule = &SimpleMappingRule{
		MappingRuleBase: MappingRuleBase{
			Output: NewRuleTargetButton("", t.outputDevice, evdev.BTN_TRIGGER, false),
			Modes:  []string{"*"},
		},
		Input: NewRuleTargetButton("", t.inputDevice, evdev.BTN_TRIGGER, false),
	}

	t.invertedRule = &SimpleMappingRule{
		MappingRuleBase: MappingRuleBase{
			Output: NewRuleTargetButton("", t.outputDevice, evdev.BTN_TRIGGER, false),
			Modes:  []string{"*"},
		},
		Input: NewRuleTargetButton("", t.inputDevice, evdev.BTN_TRIGGER, true),
	}
}

func (t *SimpleMappingRuleTests) TestMatchEvent() {
	// A matching input event should produce an output event
	correctOutput := &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  evdev.BTN_TRIGGER,
		Value: 1,
	}

	event := t.sampleRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.EqualValues(correctOutput, event)

	// An input event from the wrong device should produce a nil event
	event = t.sampleRule.MatchEvent(
		t.wrongInputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.Nil(event)

	// An input event from the wrong button should produce a nil event
	event = t.sampleRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TOP, Value: 1}, t.mode)
	t.Nil(event)
}

func (t *SimpleMappingRuleTests) TestMatchEventInverted() {
	// A matching input event should produce an output event
	correctOutput := &evdev.InputEvent{
		Type: evdev.EV_KEY,
		Code: evdev.BTN_TRIGGER,
	}

	// Should get the opposite value out that we send in
	correctOutput.Value = 0
	event := t.invertedRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.EqualValues(correctOutput, event)

	correctOutput.Value = 1
	event = t.invertedRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 0}, t.mode)
	t.EqualValues(correctOutput, event)
}

func TestRunnerMatching(t *testing.T) {
	suite.Run(t, new(SimpleMappingRuleTests))
}
