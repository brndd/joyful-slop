package config

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

type MakeRuleTargetsTests struct {
	suite.Suite
	devs map[string]*evdev.InputDevice
}

func (t *MakeRuleTargetsTests) SetupSuite() {
	t.devs = map[string]*evdev.InputDevice{
		"test": {},
	}
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetButton() {
	config := RuleTargetConfig{
		Device: "test",
	}
	t.Run("Standard keycode", func() {
		config.Button = "BTN_TRIGGER"
		rule, err := makeRuleTargetButton(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_TRIGGER, rule.Button)
	})

	t.Run("Hex code", func() {
		config.Button = "0x2fd"
		rule, err := makeRuleTargetButton(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.EvCode(0x2fd), rule.Button)
	})

	t.Run("Index", func() {
		config.Button = "3"
		rule, err := makeRuleTargetButton(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_TOP, rule.Button)
	})

	t.Run("Index too high", func() {
		config.Button = "74"
		_, err := makeRuleTargetButton(config, t.devs)
		t.NotNil(err)
	})

	t.Run("Un-prefixed keycode", func() {
		config.Button = "pinkie"
		rule, err := makeRuleTargetButton(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_PINKIE, rule.Button)
	})

	t.Run("Invalid keycode", func() {
		config.Button = "foo"
		_, err := makeRuleTargetButton(config, t.devs)
		t.NotNil(err)
	})
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetAxis() {
	config := RuleTargetConfig{
		Device: "test",
	}

	t.Run("Standard keycode", func() {
		config.Axis = "ABS_X"
		rule, err := makeRuleTargetAxis(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.ABS_X, rule.Axis)
	})

	t.Run("Hex keycode", func() {
		config.Axis = "0x01"
		rule, err := makeRuleTargetAxis(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.ABS_Y, rule.Axis)
	})

	t.Run("Un-prefixed keycode", func() {
		config.Axis = "x"
		rule, err := makeRuleTargetAxis(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.ABS_X, rule.Axis)
	})

	t.Run("Invalid keycode", func() {
		config.Axis = "foo"
		_, err := makeRuleTargetAxis(config, t.devs)
		t.NotNil(err)
	})

	t.Run("Invalid deadzone", func() {
		config.DeadzoneEnd = 100
		config.DeadzoneStart = 1000
		_, err := makeRuleTargetAxis(config, t.devs)
		t.NotNil(err)
	})
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetRelaxis() {
	config := RuleTargetConfig{
		Device: "test",
	}

	t.Run("Standard keycode", func() {
		config.Axis = "REL_WHEEL"
		rule, err := makeRuleTargetRelaxis(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_WHEEL, rule.Axis)
	})

	t.Run("Hex keycode", func() {
		config.Axis = "0x00"
		rule, err := makeRuleTargetRelaxis(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_X, rule.Axis)
	})

	t.Run("Un-prefixed keycode", func() {
		config.Axis = "wheel"
		rule, err := makeRuleTargetRelaxis(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_WHEEL, rule.Axis)
	})

	t.Run("Invalid keycode", func() {
		config.Axis = "foo"
		_, err := makeRuleTargetRelaxis(config, t.devs)
		t.NotNil(err)
	})

	t.Run("Incorrect axis type", func() {
		config.Axis = "ABS_X"
		_, err := makeRuleTargetRelaxis(config, t.devs)
		t.NotNil(err)
	})
}

func TestRunnerMakeRuleTargets(t *testing.T) {
	suite.Run(t, new(MakeRuleTargetsTests))
}
