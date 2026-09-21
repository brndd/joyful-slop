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
