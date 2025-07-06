package mappingrules

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

type MappingRuleButtonTests struct {
	suite.Suite
	inputDevice      *evdev.InputDevice
	wrongInputDevice *evdev.InputDevice
	outputDevice     *evdev.InputDevice
	mode             *string
	sampleRule       *MappingRuleButton
	invertedRule     *MappingRuleButton
}

func (t *MappingRuleButtonTests) SetupTest() {
	t.inputDevice = &evdev.InputDevice{}
	t.wrongInputDevice = &evdev.InputDevice{}
	t.outputDevice = &evdev.InputDevice{}
	mode := "*"
	t.mode = &mode

	// TODO: implement a constructor function...
	t.sampleRule = &MappingRuleButton{
		MappingRuleBase: MappingRuleBase{
			Modes: []string{"*"},
		},
		Input:  NewRuleTargetButton("", t.inputDevice, evdev.BTN_TRIGGER, false),
		Output: NewRuleTargetButton("", t.outputDevice, evdev.BTN_TRIGGER, false),
	}

	t.invertedRule = &MappingRuleButton{
		MappingRuleBase: MappingRuleBase{
			Modes: []string{"*"},
		},
		Output: NewRuleTargetButton("", t.outputDevice, evdev.BTN_TRIGGER, false),
		Input:  NewRuleTargetButton("", t.inputDevice, evdev.BTN_TRIGGER, true),
	}
}

func (t *MappingRuleButtonTests) TestMatchEvent() {
	// A matching input event should produce an output event
	correctOutput := &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  evdev.BTN_TRIGGER,
		Value: 1,
	}

	_, event := t.sampleRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.EqualValues(correctOutput, event)

	// An input event from the wrong device should produce a nil event
	_, event = t.sampleRule.MatchEvent(
		t.wrongInputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.Nil(event)

	// An input event from the wrong button should produce a nil event
	_, event = t.sampleRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TOP, Value: 1}, t.mode)
	t.Nil(event)
}

func (t *MappingRuleButtonTests) TestMatchEventInverted() {
	// A matching input event should produce an output event
	correctOutput := &evdev.InputEvent{
		Type: evdev.EV_KEY,
		Code: evdev.BTN_TRIGGER,
	}

	// Should get the opposite value out that we send in
	correctOutput.Value = 0
	_, event := t.invertedRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.EqualValues(correctOutput, event)

	correctOutput.Value = 1
	_, event = t.invertedRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 0}, t.mode)
	t.EqualValues(correctOutput, event)
}

func TestRunnerMatching(t *testing.T) {
	suite.Run(t, new(MappingRuleButtonTests))
}
