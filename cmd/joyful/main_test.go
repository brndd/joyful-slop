package main

import (
	"testing"

	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
	"github.com/stretchr/testify/require"
)

type modeTestRule struct {
	seenMode      string
	requestedMode string
	transitioned  bool
	initiated     bool
	reset         bool
	resetMode     string
	active        bool
}

type modeChangingTestRule struct {
	*modeTestRule
}

func (*modeChangingTestRule) ChangesMode() {}
func (rule *modeChangingTestRule) ModeChangeActive() bool {
	return rule.active
}

func (rule *modeTestRule) MatchEvent(_ mappingrules.Device, _ *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	rule.seenMode = *mode
	if rule.requestedMode != "" {
		*mode = rule.requestedMode
	}
	return nil, nil
}

func (rule *modeTestRule) ModeChanged(_ string, initiated bool) []mappingrules.OutputEvent {
	rule.transitioned = true
	rule.initiated = initiated
	return nil
}

func (rule *modeTestRule) Reset(mode *string) []mappingrules.OutputEvent {
	rule.reset = true
	if rule.resetMode != "" {
		*mode = rule.resetMode
	}
	return nil
}

func TestShouldAnnounceModeChange(t *testing.T) {
	require.True(t, shouldAnnounceModeChange("SCM Mode", "Nav Mode", false))
	require.False(t, shouldAnnounceModeChange("SCM Mode", "SCM Mode", false))
	require.False(t, shouldAnnounceModeChange("SCM Mode", "Modifier", true))
	require.False(t, shouldAnnounceModeChange("Modifier", "SCM Mode", true))
}

func TestMatchRulesDefersModeChange(t *testing.T) {
	modeRule := &modeChangingTestRule{modeTestRule: &modeTestRule{requestedMode: "b"}}
	followingRule := &modeTestRule{}
	rules := []mappingrules.MappingRule{modeRule, followingRule}

	_, requestedMode, initiator := matchRules(rules, nil, &evdev.InputEvent{}, "a")

	require.Equal(t, "a", modeRule.seenMode)
	require.Equal(t, "a", followingRule.seenMode)
	require.Equal(t, "b", requestedMode)
	require.Same(t, modeRule, initiator)
}

func TestMatchRulesIgnoresAdditionalModeRequests(t *testing.T) {
	firstRule := &modeChangingTestRule{modeTestRule: &modeTestRule{requestedMode: "b"}}
	secondRule := &modeChangingTestRule{modeTestRule: &modeTestRule{requestedMode: "c"}}
	rules := []mappingrules.MappingRule{firstRule, secondRule}

	_, requestedMode, initiator := matchRules(rules, nil, &evdev.InputEvent{}, "a")

	require.Equal(t, "b", requestedMode)
	require.Same(t, firstRule, initiator)
	require.Empty(t, secondRule.seenMode, "a discarded mode-changing rule must not mutate state")
}

func TestMatchRulesLetsActiveModeRuleHandleRelease(t *testing.T) {
	firstRule := &modeChangingTestRule{modeTestRule: &modeTestRule{requestedMode: "b"}}
	activeRule := &modeChangingTestRule{modeTestRule: &modeTestRule{requestedMode: "c", active: true}}
	rules := []mappingrules.MappingRule{firstRule, activeRule}

	_, requestedMode, initiator := matchRules(rules, nil, &evdev.InputEvent{}, "a")

	require.Equal(t, "b", requestedMode)
	require.Same(t, firstRule, initiator)
	require.Equal(t, "a", activeRule.seenMode, "an active rule must receive its cleanup event")
}

func TestModeTransitionIdentifiesInitiator(t *testing.T) {
	modeRule := &modeTestRule{}
	otherRule := &modeTestRule{}
	rules := []mappingrules.MappingRule{modeRule, otherRule}

	modeTransitionEvents(rules, "b", modeRule)

	require.True(t, modeRule.transitioned)
	require.True(t, modeRule.initiated)
	require.True(t, otherRule.transitioned)
	require.False(t, otherRule.initiated)
}

func TestResetRuleEventsResetsRulesAndMode(t *testing.T) {
	firstRule := &modeTestRule{}
	shiftRule := &modeTestRule{resetMode: "base"}
	rules := []mappingrules.MappingRule{firstRule, shiftRule}
	mode := "modifier"

	resetRuleEvents(rules, &mode, "startup")

	require.True(t, firstRule.reset)
	require.True(t, shiftRule.reset)
	require.Equal(t, "startup", mode)
}
