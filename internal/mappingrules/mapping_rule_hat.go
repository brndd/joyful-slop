package mappingrules

import (
	"git.annabunches.net/annabunches/joyful/internal/configparser"
	"github.com/holoplot/go-evdev"
)

// A Simple Mapping Rule can map a button to a button or an axis to an axis.
type MappingRuleHat struct {
	MappingRuleBase
	Input  *RuleTargetHat
	Output *RuleTargetHat
	active bool
}

func NewMappingRuleHat(ruleConfig configparser.RuleConfigHat,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base MappingRuleBase) (*MappingRuleHat, error) {

	input, err := NewRuleTargetHatFromConfig(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := NewRuleTargetHatFromConfig(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &MappingRuleHat{
		MappingRuleBase: base,
		Input:           input,
		Output:          output,
	}, nil
}

func (rule *MappingRuleHat) MatchEvent(device Device, event *evdev.InputEvent, mode *string) (*evdev.InputDevice, *evdev.InputEvent) {
	if !rule.MappingRuleBase.modeCheck(mode) ||
		!rule.Input.MatchEvent(device, event) {
		return nil, nil
	}

	value := rule.Input.NormalizeValue(event.Value)
	rule.active = value != 0
	// The cast here is safe because the interface is only ever different for unit tests
	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(value, mode)
}

func (rule *MappingRuleHat) ModeChanged(newMode string, initiated bool) []OutputEvent {
	if initiated || rule.MappingRuleBase.modeMatches(newMode) {
		return nil
	}
	return rule.deactivate()
}

func (rule *MappingRuleHat) Reset(_ *string) []OutputEvent {
	return rule.deactivate()
}

func (rule *MappingRuleHat) deactivate() []OutputEvent {
	if !rule.active {
		return nil
	}
	rule.active = false
	return []OutputEvent{{Device: rule.Output.Device.(*evdev.InputDevice), Event: rule.Output.CreateEvent(0, nil)}}
}
