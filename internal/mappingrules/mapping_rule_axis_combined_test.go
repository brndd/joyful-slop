package mappingrules

import (
	"fmt"
	"testing"

	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/suite"
)

// TODO: revisit all of this after adding new functionality to
// RuleTargetAxis
type MappingRuleAxisCombinedTests struct {
	suite.Suite

	inputDevice      *InputDeviceMock
	outputDevice     *evdev.InputDevice
	inputTargetLower *RuleTargetAxis
	inputTargetUpper *RuleTargetAxis
	outputTarget     *RuleTargetAxis

	base MappingRuleBase
	mode *string
}

func TestRunnerMappingRuleAxisCombined(t *testing.T) {
	suite.Run(t, new(MappingRuleAxisCombinedTests))
}

func (t *MappingRuleAxisCombinedTests) SetupTest() {
	mode := "*"
	t.mode = &mode

	t.inputDevice = NewInputDeviceMock()
	t.inputDevice.Stub("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
		evdev.ABS_X: {Minimum: 0, Maximum: 10000},
		evdev.ABS_Y: {Minimum: 0, Maximum: 10000},
	}, nil)

	t.inputTargetLower, _ = NewRuleTargetAxis("test-input", t.inputDevice, evdev.ABS_X, true, 0, 0)
	t.inputTargetUpper, _ = NewRuleTargetAxis("test-input", t.inputDevice, evdev.ABS_Y, false, 0, 0)

	t.outputDevice = &evdev.InputDevice{}
	t.outputTarget, _ = NewRuleTargetAxis("test-output", t.outputDevice, evdev.ABS_X, false, 0, 0)

	t.base = NewMappingRuleBase("", []string{"*"})

	// We clear the AbsInfo call here so it can be cleanly set by the (sub-)tests
	t.inputDevice.Reset()
}

func (t *MappingRuleAxisCombinedTests) TearDownTest() {
	t.inputDevice.Reset()
}

func (t *MappingRuleAxisCombinedTests) TearDownSubTest() {
	t.inputDevice.Reset()
}

func (t *MappingRuleAxisCombinedTests) TestNewMappingRuleAxisCombined() {
	t.inputDevice.Stub("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
		evdev.ABS_X: {Minimum: 0, Maximum: 10000},
		evdev.ABS_Y: {Minimum: 0, Maximum: 10000},
	}, nil)

	rule := NewMappingRuleAxisCombined(t.base, t.inputTargetLower, t.inputTargetUpper, t.outputTarget)
	t.EqualValues(0, rule.InputLower.OutputMax)
	t.EqualValues(0, rule.InputUpper.OutputMin)
}

func (t *MappingRuleAxisCombinedTests) TestMatchEvent() {
	rule := NewMappingRuleAxisCombined(t.base, t.inputTargetLower, t.inputTargetUpper, t.outputTarget)

	t.Run("Lower Input", func() {
		testCases := []struct{ in, out int32 }{
			{10000, AxisValueMin},
			{0, 0},
			{5000, AxisValueMin / 2},
		}

		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%d->%d", testCase.in, testCase.out), func() {
				t.inputDevice.Stub("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
					evdev.ABS_X: {Minimum: 0, Maximum: 10000, Value: testCase.in},
					evdev.ABS_Y: {Minimum: 0, Maximum: 10000, Value: 0},
				}, nil)

				device, event := rule.MatchEvent(t.inputDevice, &evdev.InputEvent{Type: evdev.EV_ABS, Code: evdev.ABS_X, Value: testCase.in}, t.mode)

				t.Equal(t.outputDevice, device)
				t.InDelta(testCase.out, event.Value, 1)
			})
		}
	})

	t.Run("Upper Input", func() {
		testCases := []struct{ in, out int32 }{
			{10000, AxisValueMax},
			{0, 0},
			{5000, AxisValueMax / 2},
		}

		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%d->%d", testCase.in, testCase.out), func() {
				t.inputDevice.Stub("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
					evdev.ABS_X: {Minimum: 0, Maximum: 10000, Value: 0},
					evdev.ABS_Y: {Minimum: 0, Maximum: 10000, Value: testCase.in},
				}, nil)

				device, event := rule.MatchEvent(t.inputDevice, &evdev.InputEvent{Type: evdev.EV_ABS, Code: evdev.ABS_Y, Value: testCase.in}, t.mode)

				t.Equal(t.outputDevice, device)
				t.InDelta(testCase.out, event.Value, 1)
			})
		}
	})

	t.Run("Combined Inputs", func() {
		testCases := []struct{ x, y, out int32 }{
			{0, 0, 0},
			{5000, 5000, 0},
			{10000, 10000, 0},
			{5000, 10000, AxisValueMax / 2},
		}

		for _, testCase := range testCases {
			t.Run(fmt.Sprintf("%d,%d->%d", testCase.x, testCase.y, testCase.out), func() {
				t.inputDevice.Stub("AbsInfos").Return(map[evdev.EvCode]evdev.AbsInfo{
					evdev.ABS_X: {Minimum: 0, Maximum: 10000, Value: testCase.x},
					evdev.ABS_Y: {Minimum: 0, Maximum: 10000, Value: testCase.y},
				}, nil)

				device, event := rule.MatchEvent(t.inputDevice, &evdev.InputEvent{Type: evdev.EV_ABS, Code: evdev.ABS_Y, Value: testCase.y}, t.mode)

				t.Equal(t.outputDevice, device)
				t.InDelta(testCase.out, event.Value, 1)
			})
		}
	})

	// TODO: add tests for exception cases
}
