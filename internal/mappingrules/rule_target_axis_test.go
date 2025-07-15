package mappingrules

import (
	"errors"
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RuleTargetAxisTests struct {
	suite.Suite
	mock *InputDeviceMock
	call *mock.Call
}

func (t *RuleTargetAxisTests) SetupTest() {
	t.mock = new(InputDeviceMock)
	t.call = t.mock.On("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
		evdev.ABS_X: {
			Minimum: 0,
			Maximum: 10000,
		},
		evdev.ABS_Y: {
			Minimum: -10000,
			Maximum: 10000,
		},
	}, nil)
}

func (t *RuleTargetAxisTests) TearDownTest() {
	t.call.Unset()
}

func (t *RuleTargetAxisTests) TestNewRuleTargetAxis() {
	// RuleTargets should get created
	ruleTarget, err := NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 0)
	t.Nil(err)
	t.EqualValues(10000, ruleTarget.axisSize)

	ruleTarget, err = NewRuleTargetAxis("", t.mock, evdev.ABS_Y, false, 0, 0)
	t.Nil(err)
	t.EqualValues(20000, ruleTarget.axisSize)

	// Creating a rule with a deadzone should work and reduce the axisSize
	ruleTarget, err = NewRuleTargetAxis("", t.mock, evdev.ABS_Y, false, -500, 500)
	t.Nil(err)
	t.EqualValues(19000, ruleTarget.axisSize)
	t.EqualValues(-500, ruleTarget.DeadzoneStart)
	t.EqualValues(500, ruleTarget.DeadzoneEnd)

	// Creating a rule with a deadzone should fail if end > start
	_, err = NewRuleTargetAxis("", t.mock, evdev.ABS_Y, false, 500, -500)
	t.NotNil(err)

	// Creating a rule on a non-existent axis should err
	_, err = NewRuleTargetAxis("", t.mock, evdev.ABS_Z, false, 0, 0)
	t.NotNil(err)

	// If Absinfo has an error, we should create a device with permissive bounds
	t.call.Unset()
	t.mock.On("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{}, errors.New("Test Error"))
	ruleTarget, err = NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 0)
	t.Nil(err)
	t.Equal(AxisValueMax-AxisValueMin, ruleTarget.axisSize)
}

func (t *RuleTargetAxisTests) TestNormalizeValue() {
	// Basic normalization should work
	ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 0)
	t.Equal(AxisValueMax, ruleTarget.NormalizeValue(int32(10000)))
	t.Equal(AxisValueMin, ruleTarget.NormalizeValue(int32(0)))
	t.EqualValues(0, ruleTarget.NormalizeValue(int32(5000)))

	// Normalization with a deadzone should work
	ruleTarget, _ = NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 5000)
	t.Equal(AxisValueMax, ruleTarget.NormalizeValue(int32(10000)))
	t.True(ruleTarget.NormalizeValue(int32(5001)) < int32(-31000))
	t.EqualValues(0, ruleTarget.NormalizeValue(int32(7500)))

	// Normalization on an inverted axis should work
	ruleTarget, _ = NewRuleTargetAxis("", t.mock, evdev.ABS_X, true, 0, 0)
	t.Equal(AxisValueMax, ruleTarget.NormalizeValue(int32(0)))
	t.Equal(AxisValueMin, ruleTarget.NormalizeValue(int32(10000)))

	// Normalization past the stated axis bounds should clamp
	ruleTarget, _ = NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 0)
	t.Equal(AxisValueMin, ruleTarget.NormalizeValue(int32(-30000)))
	t.Equal(AxisValueMax, ruleTarget.NormalizeValue(int32(30000)))
}

func (t *RuleTargetAxisTests) TestMatchEvent() {
	ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_Y, false, -500, 500)
	validEvent := &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  evdev.ABS_Y,
		Value: 800,
	}
	deadzoneEvent := &evdev.InputEvent{
		Type:  evdev.EV_ABS,
		Code:  evdev.ABS_Y,
		Value: 200,
	}

	// An event on the correct device and axis should match
	t.True(ruleTarget.MatchEvent(t.mock, validEvent))

	// A value on the wrong device should not match
	t.False(ruleTarget.MatchEvent(&evdev.InputDevice{}, validEvent))

	// A value in the deadzone should not match
	t.False(ruleTarget.MatchEvent(t.mock, deadzoneEvent))
}

func (t *RuleTargetAxisTests) TestCreateEvent() {
	ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 0)
	expected := &evdev.InputEvent{
		Type: evdev.EV_ABS,
		Code: evdev.ABS_X,
	}

	// Basic event creation
	testValue := int32(3928) // Arbitrarily chosen test value
	expected.Value = testValue
	t.EqualValues(expected, ruleTarget.CreateEvent(testValue, nil))

	// Validate axis clamping
	testValue = int32(64000)
	expected.Value = AxisValueMax
	t.EqualValues(expected, ruleTarget.CreateEvent(testValue, nil))

	testValue = int32(-64000)
	expected.Value = AxisValueMin
	t.EqualValues(expected, ruleTarget.CreateEvent(testValue, nil))
}

func (t *RuleTargetAxisTests) TestGetAxisStrength() {
	t.Run("With no deadzone", func() {
		ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 0)
		t.Equal(0.0, ruleTarget.GetAxisStrength(0))
		t.Equal(1.0, ruleTarget.GetAxisStrength(10000))
		t.Equal(0.5, ruleTarget.GetAxisStrength(5000))
	})

	t.Run("With low deadzone", func() {
		ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 0, 5000)
		t.InDelta(0.0, ruleTarget.GetAxisStrength(5001), 0.01)
		t.InDelta(0.5, ruleTarget.GetAxisStrength(7500), 0.01)
		t.Equal(1.0, ruleTarget.GetAxisStrength(10000))
	})

	t.Run("With high deadzone", func() {
		ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, false, 5000, 10000)
		t.Equal(0.0, ruleTarget.GetAxisStrength(0))
		t.InDelta(0.5, ruleTarget.GetAxisStrength(2500), 0.01)
		t.InDelta(1.0, ruleTarget.GetAxisStrength(4999), 0.01)
	})

	t.Run("Inverted", func() {
		ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, true, 0, 0)
		t.Equal(1.0, ruleTarget.GetAxisStrength(0))
		t.Equal(0.5, ruleTarget.GetAxisStrength(5000))
		t.Equal(0.0, ruleTarget.GetAxisStrength(10000))
	})

	t.Run("Inverted with low deadzone", func() {
		ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, true, 0, 5000)
		t.InDelta(1.0, ruleTarget.GetAxisStrength(5001), 0.01)
		t.InDelta(0.5, ruleTarget.GetAxisStrength(7500), 0.01)
		t.Equal(0.0, ruleTarget.GetAxisStrength(10000))
	})

	t.Run("Inverted with high deadzone", func() {
		ruleTarget, _ := NewRuleTargetAxis("", t.mock, evdev.ABS_X, true, 5000, 10000)
		t.InDelta(0.0, ruleTarget.GetAxisStrength(4999), 0.01)
		t.InDelta(0.5, ruleTarget.GetAxisStrength(2500), 0.01)
		t.Equal(1.0, ruleTarget.GetAxisStrength(0))
	})
}

func TestRunnerRuleTargetAxisTests(t *testing.T) {
	suite.Run(t, new(RuleTargetAxisTests))
}
