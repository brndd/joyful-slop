package config

import (
	"fmt"
	"strings"

	"git.annabunches.net/annabunches/joyful/internal/logger"
	"git.annabunches.net/annabunches/joyful/internal/mappingrules"
	"github.com/holoplot/go-evdev"
)

func (parser *ConfigParser) BuildRules(pDevs map[string]*evdev.InputDevice, vDevs map[string]*evdev.InputDevice) []mappingrules.MappingRule {
	rules := make([]mappingrules.MappingRule, 0)

	for _, ruleConfig := range parser.config.Rules {
		var newRule mappingrules.MappingRule
		var err error
		switch strings.ToLower(ruleConfig.Type) {
		case RuleTypeSimple:
			newRule, err = makeSimpleRule(ruleConfig, pDevs, vDevs)
		case RuleTypeCombo:
			newRule, err = makeComboRule(ruleConfig, pDevs, vDevs)
		case RuleTypeLatched:
			newRule, err = makeLatchedRule(ruleConfig, pDevs, vDevs)
		}

		if err != nil {
			logger.LogError(err, "")
			continue
		}
		rules = append(rules, newRule)
	}

	return rules
}

func makeSimpleRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, vDevs map[string]*evdev.InputDevice) (*mappingrules.SimpleMappingRule, error) {
	input, err := makeRuleTarget(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTarget(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.SimpleMappingRule{
		MappingRuleBase: mappingrules.MappingRuleBase{
			Output: output,
		},
		Input:  input,
		Name:   ruleConfig.Name,
	}, nil
}

func makeComboRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, vDevs map[string]*evdev.InputDevice) (*mappingrules.ComboMappingRule, error) {
	inputs := make([]mappingrules.RuleTarget, 0)
	for _, inputConfig := range ruleConfig.Inputs {
		input, err := makeRuleTarget(inputConfig, pDevs)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, input)
	}

	output, err := makeRuleTarget(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.ComboMappingRule{
		MappingRuleBase: mappingrules.MappingRuleBase{
			Output: output,
		},
		Inputs: inputs,
		State:   0,
		Name:   ruleConfig.Name,
	}, nil
}

func makeLatchedRule(ruleConfig RuleConfig, pDevs map[string]*evdev.InputDevice, vDevs map[string]*evdev.InputDevice) (*mappingrules.LatchedMappingRule, error) {
	input, err := makeRuleTarget(ruleConfig.Input, pDevs)
	if err != nil {
		return nil, err
	}

	output, err := makeRuleTarget(ruleConfig.Output, vDevs)
	if err != nil {
		return nil, err
	}

	return &mappingrules.LatchedMappingRule{
		MappingRuleBase: mappingrules.MappingRuleBase{
			Output: output,
		},
		Input:  input,
		Name:   ruleConfig.Name,
		State:  false,
	}, nil
}

// makeInputRuleTarget takes an Input declaration from the YAML and returns a fully formed RuleTarget.
func makeRuleTarget(targetConfig RuleTargetConfig, devs map[string]*evdev.InputDevice) (mappingrules.RuleTarget, error) {
	ruleTarget := mappingrules.RuleTarget{}

	device, ok := devs[targetConfig.Device]
	if !ok {
		return mappingrules.RuleTarget{}, fmt.Errorf("couldn't build rule due to non-existent device '%s'", targetConfig.Device)
	}
	ruleTarget.Device = device

	eventType, eventCode, err := decodeRuleTargetValues(targetConfig)
	if err != nil {
		return ruleTarget, err
	}
	ruleTarget.Type = eventType
	ruleTarget.Code = eventCode
	ruleTarget.Inverted = targetConfig.Inverted
	ruleTarget.DeviceName = targetConfig.Device

	return ruleTarget, nil
}

// decodeRuleTargetValues returns the appropriate evdev.EvType and evdev.EvCode values
// for a given RuleTargetConfig, converting the config file strings into appropriate constants
//
// Todo: support different formats for key specification
func decodeRuleTargetValues(target RuleTargetConfig) (evdev.EvType, evdev.EvCode, error) {
	var eventType evdev.EvType
	var eventCode evdev.EvCode
	var ok bool

	if target.Button != "" {
		eventType = evdev.EV_KEY
		eventCode, ok = evdev.KEYFromString[target.Button]
		if !ok {
			return 0, 0, fmt.Errorf("skipping rule due to invalid button code '%s'", target.Button)
		}
	} else if target.Axis != "" {
		eventType = evdev.EV_ABS
		eventCode, ok = evdev.ABSFromString[target.Axis]
		if !ok {
			return 0, 0, fmt.Errorf("skipping rule due to invalid axis code '%s'", target.Button)
		}
	}

	return eventType, eventCode, nil
}
