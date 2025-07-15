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
	base             MappingRuleBase
}

func (t *MappingRuleButtonTests) SetupTest() {
	t.inputDevice = &evdev.InputDevice{}
	t.wrongInputDevice = &evdev.InputDevice{}
	t.outputDevice = &evdev.InputDevice{}
	mode := "*"
	t.mode = &mode
	t.base = NewMappingRuleBase("", []string{})
}

func (t *MappingRuleButtonTests) TestMatchEvent() {
	inputButton, _ := NewRuleTargetButton("", t.inputDevice, evdev.BTN_TRIGGER, false)
	outputButton, _ := NewRuleTargetButton("", t.outputDevice, evdev.BTN_TRIGGER, false)
	testRule := NewMappingRuleButton(t.base, inputButton, outputButton)

	// A matching input event should produce an output event
	expected := &evdev.InputEvent{
		Type:  evdev.EV_KEY,
		Code:  evdev.BTN_TRIGGER,
		Value: 1,
	}

	_, event := testRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.EqualValues(expected, event)

	// An input event from the wrong device should produce a nil event
	_, event = testRule.MatchEvent(
		t.wrongInputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.Nil(event)

	// An input event from the wrong button should produce a nil event
	_, event = testRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TOP, Value: 1}, t.mode)
	t.Nil(event)
}

func (t *MappingRuleButtonTests) TestMatchEventInverted() {
	inputButton, _ := NewRuleTargetButton("", t.inputDevice, evdev.BTN_TRIGGER, true)
	outputButton, _ := NewRuleTargetButton("", t.outputDevice, evdev.BTN_TRIGGER, false)
	testRule := NewMappingRuleButton(t.base, inputButton, outputButton)

	// A matching input event should produce an output event
	expected := &evdev.InputEvent{
		Type: evdev.EV_KEY,
		Code: evdev.BTN_TRIGGER,
	}

	// Should get the opposite value out that we send in
	expected.Value = 0
	_, event := testRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 1}, t.mode)
	t.EqualValues(expected, event)

	expected.Value = 1
	_, event = testRule.MatchEvent(
		t.inputDevice,
		&evdev.InputEvent{Code: evdev.BTN_TRIGGER, Value: 0}, t.mode)
	t.EqualValues(expected, event)
}

func TestRunnerMappingRuleButtonTests(t *testing.T) {
	suite.Run(t, new(MappingRuleButtonTests))
}
