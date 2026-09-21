package mappingrules

import (
	"errors"
	"time"

	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"git.annabunches.net/annabunches/joyful/internal/logger"
	"github.com/holoplot/go-evdev"
	"github.com/jonboulle/clockwork"
)

type buttonTempoBranch struct {
	outputs []*RuleTargetButton
	mode    string
}

type MappingRuleButtonTempo struct {
	MappingRuleBase
	Input     *RuleTargetButton
	Threshold time.Duration
	tap       buttonTempoBranch
	hold      buttonTempoBranch
	clock     clockwork.Clock
	pressedAt time.Time
	active    bool
	held      bool
}

func NewMappingRuleButtonTempo(ruleConfig configparser.RuleConfigButtonTempo, pDevs, vDevs map[string]Device, modes []string, base MappingRuleBase) (*MappingRuleButtonTempo, error) {
	if ruleConfig.ThresholdMs <= 0 {
		return nil, errors.New("button-tempo threshold_ms must be positive")
	}
	input, err := NewRuleTargetButtonFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}
	tap, err := newButtonTempoBranch(ruleConfig.Tap, vDevs, modes)
	if err != nil {
		return nil, err
	}
	hold, err := newButtonTempoBranch(ruleConfig.Hold, vDevs, modes)
	if err != nil {
		return nil, err
	}
	return &MappingRuleButtonTempo{MappingRuleBase: base, Input: input, Threshold: time.Duration(ruleConfig.ThresholdMs) * time.Millisecond, tap: tap, hold: hold, clock: clockwork.NewRealClock()}, nil
}

func newButtonTempoBranch(config configparser.RuleConfigButtonTempoBranch, vDevs map[string]Device, modes []string) (buttonTempoBranch, error) {
	if config.Mode != "" && !validateModes([]string{config.Mode}, modes) {
		return buttonTempoBranch{}, errors.New("button-tempo branch specifies undefined mode")
	}
	branch := buttonTempoBranch{mode: config.Mode}
	for _, outputConfig := range config.Outputs {
		output, err := NewRuleTargetButtonFromConfig(outputConfig, vDevs)
		if err != nil {
			return buttonTempoBranch{}, err
		}
		branch.outputs = append(branch.outputs, output)
	}
	return branch, nil
}

// MatchEvent preserves the MappingRule API. Runtime callers should use MatchEvents.
func (rule *MappingRuleButtonTempo) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	events := rule.MatchEvents(device, event, mode)
	if len(events) == 0 {
		return nil, nil
	}
	return events[0].Device, events[0].Event
}

func (rule *MappingRuleButtonTempo) MatchEvents(device Device, event *evdev.InputEvent, mode *string) []OutputEvent {
	if !rule.Input.MatchEvent(device, event) {
		return nil
	}
	value := rule.Input.NormalizeValue(event.Value)
	if value != 0 {
		if rule.active || !rule.MappingRuleBase.modeCheck(mode) {
			return nil
		}
		rule.active = true
		rule.held = false
		rule.pressedAt = rule.clock.Now()
		return nil
	}
	if !rule.active {
		return nil
	}
	rule.active = false
	if rule.held {
		rule.held = false
		return branchEvents(rule.hold, 0, mode)
	}
	if !rule.clock.Now().Before(rule.pressedAt.Add(rule.Threshold)) {
		events := branchEvents(rule.hold, 1, mode)
		return append(events, branchEvents(rule.hold, 0, mode)...)
	}
	events := branchEvents(rule.tap, 1, mode)
	return append(events, branchEvents(rule.tap, 0, mode)...)
}

func (rule *MappingRuleButtonTempo) TimerEvents(mode *string) []OutputEvent {
	if !rule.active || rule.held || rule.clock.Now().Before(rule.pressedAt.Add(rule.Threshold)) {
		return nil
	}
	rule.held = true
	return branchEvents(rule.hold, 1, mode)
}

func branchEvents(branch buttonTempoBranch, value int32, mode *string) []OutputEvent {
	events := make([]OutputEvent, 0, len(branch.outputs))
	for _, output := range branch.outputs {
		events = append(events, OutputEvent{Device: output.Device.(*evdev.InputDevice), Event: output.CreateEvent(value, nil)})
	}
	if value == 1 && branch.mode != "" {
		*mode = branch.mode
		logger.Logf("Mode changed to '%s'", *mode)
	}
	return events
}
