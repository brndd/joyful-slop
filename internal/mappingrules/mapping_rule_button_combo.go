package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

// A Combo Mapping Rule can require multiple physical button presses for a single output button
type MappingRuleButtonCombo struct {
	MappingRuleBase
	Inputs []*RuleTargetButton
	Output *RuleTargetButton
	State  int
}

func NewMappingRuleButtonCombo(ruleConfig configparser.RuleConfigButtonCombo,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleButtonCombo, error) {

	inputs := make([]*RuleTargetButton, 0)
	for _, inputConfig := range ruleConfig.Inputs {
		input, err := NewRuleTargetButtonFromConfig(inputConfig, pDevs)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}

	output, err := NewRuleTargetButtonFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &MappingRuleButtonCombo{
		MappingRuleBase: base,
		Inputs:          inputs,
		Output:          output,
		State:           0,
	}, nil
}

func (rule *MappingRuleButtonCombo) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) {
		return nil, nil
	}

	// Check each of the inputs, and if we find a match, proceed
	var match *RuleTargetButton
	for _, input := range rule.Inputs {
		if input.MatchEvent(device, event) {
			match = input
			break
		}
	}

	if match == nil {
		return nil, nil
	}

	// Get the value and add/subtract it from State
	inputValue := match.NormalizeValue(event.Value)
	oldState := rule.State
	if inputValue == 0 {
		rule.State = max(rule.State-1, 0)
	}
	if inputValue == 1 {
		rule.State++
	}
	targetState := len(rule.Inputs)

	if oldState == targetState-1 && rule.State == targetState {
		return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(1, mode)
	}
	if oldState == targetState && rule.State == targetState-1 {
		return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(0, mode)
	}
	return nil, nil
}

func (rule *MappingRuleButtonCombo) ModeChanged(newMode string, initiated bool) []OutputEvent {
	if initiated || rule.MappingRuleBase.modeMatches(newMode) {
		return nil
	}
	return rule.deactivate()
}

func (rule *MappingRuleButtonCombo) Reset(_ *string) []OutputEvent {
	return rule.deactivate()
}

func (rule *MappingRuleButtonCombo) deactivate() []OutputEvent {
	wasPressed := rule.State == len(rule.Inputs)
	rule.State = 0
	if !wasPressed {
		return nil
	}
	return []OutputEvent{{Device: rule.Output.Device.(*evdev.InputDevice), Event: rule.Output.CreateEvent(0, nil)}}
}
