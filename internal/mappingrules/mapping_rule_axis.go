package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleAxis struct {
	MappingRuleBase
	Input  *RuleTargetAxis
	Output *RuleTargetAxis
	active bool
}

func NewMappingRuleAxis(ruleConfig configparser.RuleConfigAxis,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleAxis, error) {

	input, err := NewRuleTargetAxisFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetAxisFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &MappingRuleAxis{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}, nil
}

func (rule *MappingRuleAxis) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEvent(device, event) {
		return nil, nil
	}

	value := rule.Input.NormalizeValue(event.Value)
	rule.active = value != 0
	// The cast here is safe because the interface is only ever different for unit tests
	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(value, mode)
}

func (rule *MappingRuleAxis) ModeChanged(newMode string, initiated bool) []OutputEvent {
	if initiated || rule.MappingRuleBase.modeMatches(newMode) {
		return nil
	}
	return rule.deactivate()
}

func (rule *MappingRuleAxis) Reset(_ *string) []OutputEvent {
	return rule.deactivate()
}

func (rule *MappingRuleAxis) deactivate() []OutputEvent {
	if !rule.active {
		return nil
	}
	rule.active = false
	return []OutputEvent{{Device: rule.Output.Device.(*evdev.InputDevice), Event: rule.Output.CreateEvent(0, nil)}}
}
