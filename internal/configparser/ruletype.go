package configparser

import (
	"fmt"
	"strings"
)

// TODO: maybe these want to live somewhere other than configparser?
type RuleType string

const (
	RuleTypeNone          RuleType = ""
	RuleTypeButton        RuleType = "button"
	RuleTypeButtonCombo   RuleType = "button-combo"
	RuleTypeButtonLatched RuleType = "button-latched"
	RuleTypeButtonTempo   RuleType = "button-tempo"
	RuleTypeAxis          RuleType = "axis"
	RuleTypeAxisCombined  RuleType = "axis-combined"
	RuleTypeAxisToButton  RuleType = "axis-to-button"
	RuleTypeAxisToRelaxis RuleType = "axis-to-relaxis"
	RuleTypeModeSelect    RuleType = "mode-select"
	RuleTypeModeShift     RuleType = "mode-shift"
	RuleTypeHat           RuleType = "hat"
)

var (
	ruleTypeMap = map[string]RuleType{
		"button":          RuleTypeButton,
		"button-combo":    RuleTypeButtonCombo,
		"button-latched":  RuleTypeButtonLatched,
		"button-tempo":    RuleTypeButtonTempo,
		"axis":            RuleTypeAxis,
		"axis-combined":   RuleTypeAxisCombined,
		"axis-to-button":  RuleTypeAxisToButton,
		"axis-to-relaxis": RuleTypeAxisToRelaxis,
		"mode-select":     RuleTypeModeSelect,
		"mode-shift":      RuleTypeModeShift,
		"hat":             RuleTypeHat,
	}
)

func ParseRuleType(in string) (RuleType, error) {
	ruleType, ok := ruleTypeMap[strings.ToLower(in)]
	if !ok {
		return RuleTypeNone, fmt.Errorf("invalid rule type '%s'", in)
	}
	return ruleType, nil
}

func (rt *RuleType) UnmarshalYAML(unmarshal func(data interface{}) error) error {
	var raw string
	err := unmarshal(&raw)
	if err != nil {
		return err
	}

	*rt, err = ParseRuleType(raw)
	return err
}
