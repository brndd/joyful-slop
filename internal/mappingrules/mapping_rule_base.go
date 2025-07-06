package mappingrules

import "slices"

type MappingRuleBase struct {
	Name   string
	Output RuleTarget
	Modes  []string
}

func (rule *MappingRuleBase) OutputName() string {
	return rule.Output.GetDeviceName()
}

func (rule *MappingRuleBase) modeCheck(mode *string) bool {
	if rule.Modes[0] == "*" {
		return true
	}
	return slices.Contains(rule.Modes, *mode)
}
