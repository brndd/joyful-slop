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

	// The cast here is safe because the interface is only ever different for unit tests
	return rule.Output.Device.(*evdev.InputDevice), rule.Output.CreateEvent(rule.Input.NormalizeValue(event.Value), mode)
}
