package mappingrules

import (
	"testing"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/require"
)

func TestModeShiftRestoresModeActiveAtPress(t *testing.T) {
	inputDevice := new(InputDeviceMock)
	rule, err := NewMappingRuleModeShift(configparser.RuleConfigModeShift{
		Input: configparser.RuleTargetConfigButton{Device: "input", Button: "BTN_TRIGGER"},
		Mode:  "shift",
	}, map[string]Device{"input": inputDevice}, []string{"base", "other", "shift"}, NewMappingRuleBase("", nil))
	require.NoError(t, err)
	mode := "other"

	rule.MatchEvent(inputDevice, buttonInput(1), &mode)
	require.Equal(t, "shift", mode)
	rule.MatchEvent(inputDevice, &evdev.InputEvent{Type: evdev.EV_KEY, Code: evdev.BTN_TRIGGER, Value: 2}, &mode)
	require.Equal(t, "shift", mode)
	rule.MatchEvent(inputDevice, buttonInput(0), &mode)
	require.Equal(t, "other", mode)
}

func TestModeShiftResetRestoresPreviousMode(t *testing.T) {
	inputDevice := &evdev.InputDevice{}
	rule, err := NewMappingRuleModeShift(configparser.RuleConfigModeShift{
		Input: configparser.RuleTargetConfigButton{Device: "input", Button: "BTN_TRIGGER"},
		Mode:  "shift",
	}, map[string]Device{"input": inputDevice}, []string{"base", "shift"}, NewMappingRuleBase("", nil))
	require.NoError(t, err)
	mode := "base"

	rule.MatchEvent(inputDevice, buttonInput(1), &mode)
	require.Equal(t, "shift", mode)
	rule.Reset(&mode)
	require.Equal(t, "base", mode)
	require.False(t, rule.active)
}

func TestNestedModeShiftsRestoreInReverseOrder(t *testing.T) {
	firstInput := &evdev.InputDevice{}
	secondInput := &evdev.InputDevice{}
	first, err := NewMappingRuleModeShift(configparser.RuleConfigModeShift{
		Input: configparser.RuleTargetConfigButton{Device: "first", Button: "BTN_TRIGGER"},
		Mode:  "b",
	}, map[string]Device{"first": firstInput}, []string{"a", "b", "c"}, NewMappingRuleBase("", nil))
	require.NoError(t, err)
	second, err := NewMappingRuleModeShift(configparser.RuleConfigModeShift{
		Input: configparser.RuleTargetConfigButton{Device: "second", Button: "BTN_TRIGGER"},
		Mode:  "c",
	}, map[string]Device{"second": secondInput}, []string{"a", "b", "c"}, NewMappingRuleBase("", nil))
	require.NoError(t, err)
	mode := "a"

	first.MatchEvent(firstInput, buttonInput(1), &mode)
	require.Equal(t, "b", mode)
	second.MatchEvent(secondInput, buttonInput(1), &mode)
	require.Equal(t, "c", mode)
	second.MatchEvent(secondInput, buttonInput(0), &mode)
	require.Equal(t, "b", mode)
	first.MatchEvent(firstInput, buttonInput(0), &mode)
	require.Equal(t, "a", mode)
}
