package config

import (
	"fmt"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

// TODO: At some point it would *very likely* make sense to map each rule to all of the physical devices that can
// trigger it, and return that instead. Something like a map[Device][]mappingrule.MappingRule.
// This would speed up rule matching by only checking relevant rules for a given input event.
// We could take this further and make it a map[<struct of *inputdevice, type, and code>][]rule
// For very large rule-bases this may be helpful for staying performant.
func (parser *ConfigParser) BuildRules(pInputDevs map[string]*evdev.InputDevice, vInputDevs map[string]*evdev.InputDevice) []mappingrules.MappingRule {
	rules := make([]mappingrules.MappingRule, 0)
	modes := parser.GetModes()

	// Golang can't inspect the concrete map type to determine interface conformance,
	// so we handle that here.
	pDevs := make(map[string]Device)
	for name, dev := range pInputDevs {
		pDevs[name] = dev
	}
	vDevs := make(map[string]Device)
	for name, dev := range vInputDevs {
		vDevs[name] = dev
	}

	for _, ruleConfig := range parser.config.Rules {
		var newRule mappingrules.MappingRule
		var err error

		if ok := validateModes(ruleConfig.Modes, modes); !ok {
			logger.Logf("Skipping rule '%s', mode list specifies undefined mode.", ruleConfig.Name)
			continue
		}

		base := mappingrules.NewMappingRuleBase(ruleConfig.Name, ruleConfig.Modes)

		switch strings.ToLower(ruleConfig.Type) {
		case RuleTypeButton:
			newRule, err = makeMappingRuleButton(ruleConfig, pDevs, vDevs, base)
		case RuleTypeButtonCombo:
			newRule, err = makeMappingRuleCombo(ruleConfig, pDevs, vDevs, base)
		case RuleTypeLatched:
			newRule, err = makeMappingRuleLatched(ruleConfig, pDevs, vDevs, base)
		case RuleTypeAxis:
			newRule, err = makeMappingRuleAxis(ruleConfig, pDevs, vDevs, base)
		case RuleTypeAxisCombined:
			newRule, err = makeMappingRuleAxisCombined(ruleConfig, pDevs, vDevs, base)
		case RuleTypeAxisToButton:
			newRule, err = makeMappingRuleAxisToButton(ruleConfig, pDevs, vDevs, base)
		case RuleTypeAxisToRelaxis:
			newRule, err = makeMappingRuleAxisToRelaxis(ruleConfig, pDevs, vDevs, base)
		case RuleTypeModeSelect:
			newRule, err = makeMappingRuleModeSelect(ruleConfig, pDevs, modes, base)
		default:
			err = fmt.Errorf("bad rule type '%s' for rule '%s'", ruleConfig.Type, ruleConfig.Name)
		}

		if err != nil {
			logger.LogErrorf(err, "Failed to build rule '%s'", ruleConfig.Name)
			continue
		}

		rules = append(rules, newRule)
	}

	return rules
}

func makeMappingRuleButton(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleButton, error) {

	input, err := makeRuleTargetButton(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetButton(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleButton(base, input, output), nil
}

func makeMappingRuleCombo(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleButtonCombo, error) {

	inputs := make([]*mappingrules.RuleTargetButton, 0)
	for _, inputConfig := range ruleConfig.Inputs {
		input, err := makeRuleTargetButton(inputConfig, pDevs)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}

	output, err := makeRuleTargetButton(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleButtonCombo(base, inputs, output), nil
}

func makeMappingRuleLatched(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleButtonLatched, error) {

	input, err := makeRuleTargetButton(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetButton(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleButtonLatched(base, input, output), nil
}

func makeMappingRuleAxis(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleAxis, error) {

	input, err := makeRuleTargetAxis(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetAxis(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleAxis(base, input, output), nil
}

func makeMappingRuleAxisCombined(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleAxisCombined, error) {

	inputLower, err := makeRuleTargetAxis(ruleConfig.InputLower, pDevs)
	if err != nil {
		return nil, err
	}

	inputUpper, err := makeRuleTargetAxis(ruleConfig.InputUpper, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetAxis(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleAxisCombined(base, inputLower, inputUpper, output), nil
}

func makeMappingRuleAxisToButton(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleAxisToButton, error) {

	input, err := makeRuleTargetAxis(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetButton(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleAxisToButton(base, input, output, ruleConfig.RepeatRateMin, ruleConfig.RepeatRateMax), nil
}

func makeMappingRuleAxisToRelaxis(ruleConfig RuleConfig,
	pDevs map[string]Device,
	vDevs map[string]Device,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleAxisToRelaxis, error) {

	input, err := makeRuleTargetAxis(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetRelaxis(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleAxisToRelaxis(base,
		input, output,
		ruleConfig.RepeatRateMin,
		ruleConfig.RepeatRateMax,
		ruleConfig.Increment), nil
}

func makeMappingRuleModeSelect(ruleConfig RuleConfig,
	pDevs map[string]Device,
	modes []string,
	base mappingrules.MappingRuleBase) (*mappingrules.MappingRuleModeSelect, error) {

	input, err := makeRuleTargetButton(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTargetModeSelect(ruleConfig.Output, modes)
	if err != nil {
		return nil, err
	}

	return mappingrules.NewMappingRuleModeSelect(base, input, output), nil
}
