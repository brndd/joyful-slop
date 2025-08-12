// TODO: these tests should live with their rule_target_* counterparts

package mappingrules

import (
	"fmt"
	"testing"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MakeRuleTargetsTests struct {
	suite.Suite
	devs       map[string]Device
	deviceMock *DeviceMock
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

func (t *MakeRuleTargetsTests) TestMakeRuleTargetButton() {
	config := configparser.RuleTargetConfigButton{Device: "test"}

	t.Run("Standard keycode", func() {
		config.Button = "BTN_TRIGGER"
		rule, err := NewRuleTargetButtonFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_TRIGGER, rule.Button)
	})

	t.Run("Hex code", func() {
		config.Button = "0x2fd"
		rule, err := NewRuleTargetButtonFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.EvCode(0x2fd), rule.Button)
	})

	t.Run("Index", func() {
		config.Button = "3"
		rule, err := NewRuleTargetButtonFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_TOP, rule.Button)
	})

	t.Run("Index too high", func() {
		config.Button = "74"
		_, err := NewRuleTargetButtonFromConfig(config, t.devs)
		t.NotNil(err)
	})

	t.Run("Un-prefixed keycode", func() {
		config.Button = "pinkie"
		rule, err := NewRuleTargetButtonFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.BTN_PINKIE, rule.Button)
	})

	t.Run("Invalid keycode", func() {
		config.Button = "foo"
		_, err := NewRuleTargetButtonFromConfig(config, t.devs)
		t.NotNil(err)
	})
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetAxis() {
	codeTestCases := []struct {
		input  string
		output evdev.EvCode
	}{
		{"ABS_X", evdev.ABS_X},
		{"0x01", evdev.ABS_Y},
		{"x", evdev.ABS_X},
	}

	for _, tc := range codeTestCases {
		t.Run(fmt.Sprintf("KeyCode %s", tc.input), func() {
			config := configparser.RuleTargetConfigAxis{Device: "test"}
			config.Axis = tc.input
			rule, err := NewRuleTargetAxisFromConfig(config, t.devs)
			t.Nil(err)
			t.EqualValues(tc.output, rule.Axis)

		})
	}

	t.Run("Invalid code", func() {
		config := configparser.RuleTargetConfigAxis{Device: "test"}
		config.Axis = "foo"
		_, err := NewRuleTargetAxisFromConfig(config, t.devs)
		t.NotNil(err)
	})

	t.Run("Invalid deadzone", func() {
		config := configparser.RuleTargetConfigAxis{Device: "test"}
		config.Axis = "x"
		config.DeadzoneEnd = 100
		config.DeadzoneStart = 1000
		_, err := NewRuleTargetAxisFromConfig(config, t.devs)
		t.NotNil(err)
	})

	relDeadzoneTestCases := []struct {
		inCenter int32
		inSize   int32
		outStart int32
		outEnd   int32
	}{
		{5000, 1000, 4500, 5500},
		{0, 500, 0, 500},
		{10000, 500, 9500, 10000},
	}

	for _, tc := range relDeadzoneTestCases {
		t.Run(fmt.Sprintf("Relative Deadzone %d +- %d", tc.inCenter, tc.inSize), func() {
			config := configparser.RuleTargetConfigAxis{
				Device:         "test",
				Axis:           "x",
				DeadzoneCenter: tc.inCenter,
				DeadzoneSize:   tc.inSize,
			}
			rule, err := NewRuleTargetAxisFromConfig(config, t.devs)

			t.Nil(err)
			t.Equal(tc.outStart, rule.DeadzoneStart)
			t.Equal(tc.outEnd, rule.DeadzoneEnd)
		})
	}

	t.Run("Deadzone center/size invalid center", func() {
		config := configparser.RuleTargetConfigAxis{
			Device:         "test",
			Axis:           "x",
			DeadzoneCenter: 20000,
			DeadzoneSize:   500,
		}
		_, err := NewRuleTargetAxisFromConfig(config, t.devs)
		t.NotNil(err)
	})

	relDeadzonePercentTestCases := []struct {
		inCenter      int32
		inSizePercent int32
		outStart      int32
		outEnd        int32
	}{
		{5000, 10, 4500, 5500},
		{0, 10, 0, 1000},
		{10000, 10, 9000, 10000},
	}

	for _, tc := range relDeadzonePercentTestCases {
		t.Run(fmt.Sprintf("Relative percent deadzone %d +- %d%%", tc.inCenter, tc.inSizePercent), func() {
			config := configparser.RuleTargetConfigAxis{
				Device:              "test",
				Axis:                "x",
				DeadzoneCenter:      tc.inCenter,
				DeadzoneSizePercent: tc.inSizePercent,
			}
			rule, err := NewRuleTargetAxisFromConfig(config, t.devs)

			t.Nil(err)
			t.Equal(tc.outStart, rule.DeadzoneStart)
			t.Equal(tc.outEnd, rule.DeadzoneEnd)
		})
	}

	t.Run("Deadzone center/percent invalid center", func() {
		config := configparser.RuleTargetConfigAxis{
			Device:              "test",
			Axis:                "x",
			DeadzoneCenter:      20000,
			DeadzoneSizePercent: 10,
		}
		_, err := NewRuleTargetAxisFromConfig(config, t.devs)
		t.NotNil(err)
	})
}

func (t *MakeRuleTargetsTests) TestMakeRuleTargetRelaxis() {
	config := configparser.RuleTargetConfigRelaxis{Device: "test"}

	t.Run("Standard keycode", func() {
		config.Axis = "REL_WHEEL"
		rule, err := NewRuleTargetRelaxisFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_WHEEL, rule.Axis)
	})

	t.Run("Hex keycode", func() {
		config.Axis = "0x00"
		rule, err := NewRuleTargetRelaxisFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_X, rule.Axis)
	})

	t.Run("Un-prefixed keycode", func() {
		config.Axis = "wheel"
		rule, err := NewRuleTargetRelaxisFromConfig(config, t.devs)
		t.Nil(err)
		t.EqualValues(evdev.REL_WHEEL, rule.Axis)
	})

	t.Run("Invalid keycode", func() {
		config.Axis = "foo"
		_, err := NewRuleTargetRelaxisFromConfig(config, t.devs)
		t.NotNil(err)
	})

	t.Run("Incorrect axis type", func() {
		config.Axis = "ABS_X"
		_, err := NewRuleTargetRelaxisFromConfig(config, t.devs)
		t.NotNil(err)
	})
}
