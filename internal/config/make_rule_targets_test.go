package config

import (
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MakeRuleTargetsTests struct {
	suite.Suite
	devs       map[string]Device
	deviceMock *DeviceMock
	config     RuleTargetConfig
}

type DeviceMock struct {
	mock.Mock
}

func (m *DeviceMock) AbsInfos() (map[evdev.EvCode]evdev.AbsInfo, error) {
	args := m.Called()
	return args.Get(0).(map[evdev.EvCode]evdev.AbsInfo), args.Error(1)
}

func TestRunnerMakeRuleTargets(t *testing.T) {
	suite.Run(t, new(MakeRuleTargetsTests))
}

func (t *MakeRuleTargetsTests) SetupSuite() {
	t.deviceMock = new(DeviceMock)
	t.deviceMock.On("AbsInfos").Return(
		map[evdev.EvCode]evdev.AbsInfo{
			evdev.ABS_X: {
				Minimum: 0,
				Maximum: 10000,
			},
			evdev.ABS_Y: {
				Minimum: 0,
				Maximum: 10000,
			},
		}, nil,
	)
	t.devs = map[string]Device{
		"test": t.deviceMock,
	}
}

func (t *MakeRuleTargetsTests) SetupSubTest() {
	t.config = RuleTargetConfig{
		Device: "test",
	}
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetButton() {
	t.Run("Standard keycode", func() {
		t.config.Button = "BTN_TRIGGER"
		rule, err := makeRuleTargetButton(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_TRIGGER, rule.Button)
	})

	t.Run("Hex code", func() {
		t.config.Button = "0x2fd"
		rule, err := makeRuleTargetButton(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.EvCode(0x2fd), rule.Button)
	})

	t.Run("Index", func() {
		t.config.Button = "3"
		rule, err := makeRuleTargetButton(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_TOP, rule.Button)
	})

	t.Run("Index too high", func() {
		t.config.Button = "74"
		_, err := makeRuleTargetButton(t.config, t.devs)
		t.NotNil(err)
	})

	t.Run("Un-prefixed keycode", func() {
		t.config.Button = "pinkie"
		rule, err := makeRuleTargetButton(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_PINKIE, rule.Button)
	})

	t.Run("Invalid keycode", func() {
		t.config.Button = "foo"
		_, err := makeRuleTargetButton(t.config, t.devs)
		t.NotNil(err)
	})
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetAxis() {
	t.Run("Standard code", func() {
		t.config.Axis = "ABS_X"
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.ABS_X, rule.Axis)
	})

	t.Run("Hex code", func() {
		t.config.Axis = "0x01"
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.ABS_Y, rule.Axis)
	})

	t.Run("Un-prefixed code", func() {
		t.config.Axis = "x"
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.ABS_X, rule.Axis)
	})

	t.Run("Invalid code", func() {
		t.config.Axis = "foo"
		_, err := makeRuleTargetAxis(t.config, t.devs)
		t.NotNil(err)
	})

	t.Run("Invalid deadzone", func() {
		t.config.Axis = "x"
		t.config.DeadzoneEnd = 100
		t.config.DeadzoneStart = 1000
		_, err := makeRuleTargetAxis(t.config, t.devs)
		t.NotNil(err)
	})

	t.Run("Deadzone center/size", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 5000
		t.config.DeadzoneSize = 1000
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(4500, rule.DeadzoneStart)
		t.EqualValues(5500, rule.DeadzoneEnd)
	})

	t.Run("Deadzone center/size lower boundary", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 0
		t.config.DeadzoneSize = 500
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(0, rule.DeadzoneStart)
		t.EqualValues(500, rule.DeadzoneEnd)
	})

	t.Run("Deadzone center/size upper boundary", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 10000
		t.config.DeadzoneSize = 500
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(9500, rule.DeadzoneStart)
		t.EqualValues(10000, rule.DeadzoneEnd)
	})

	t.Run("Deadzone center/size invalid center", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 20000
		t.config.DeadzoneSize = 500
		_, err := makeRuleTargetAxis(t.config, t.devs)
		t.NotNil(err)
	})

	t.Run("Deadzone center/percent", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 5000
		t.config.DeadzoneSizePercent = 10
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(4500, rule.DeadzoneStart)
		t.EqualValues(5500, rule.DeadzoneEnd)
	})

	t.Run("Deadzone center/percent lower boundary", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 0
		t.config.DeadzoneSizePercent = 10
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(0, rule.DeadzoneStart)
		t.EqualValues(1000, rule.DeadzoneEnd)
	})

	t.Run("Deadzone center/percent upper boundary", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 10000
		t.config.DeadzoneSizePercent = 10
		rule, err := makeRuleTargetAxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(9000, rule.DeadzoneStart)
		t.EqualValues(10000, rule.DeadzoneEnd)
	})

	t.Run("Deadzone center/percent invalid center", func() {
		t.config.Axis = "x"
		t.config.DeadzoneCenter = 20000
		t.config.DeadzoneSizePercent = 10
		_, err := makeRuleTargetAxis(t.config, t.devs)
		t.NotNil(err)
	})
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetRelaxis() {
	t.Run("Standard keycode", func() {
		t.config.Axis = "REL_WHEEL"
		rule, err := makeRuleTargetRelaxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_WHEEL, rule.Axis)
	})

	t.Run("Hex keycode", func() {
		t.config.Axis = "0x00"
		rule, err := makeRuleTargetRelaxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_X, rule.Axis)
	})

	t.Run("Un-prefixed keycode", func() {
		t.config.Axis = "wheel"
		rule, err := makeRuleTargetRelaxis(t.config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_WHEEL, rule.Axis)
	})

	t.Run("Invalid keycode", func() {
		t.config.Axis = "foo"
		_, err := makeRuleTargetRelaxis(t.config, t.devs)
		t.NotNil(err)
	})

	t.Run("Incorrect axis type", func() {
		t.config.Axis = "ABS_X"
		_, err := makeRuleTargetRelaxis(t.config, t.devs)
		t.NotNil(err)
	})
}
